use std::{
    sync::{
        atomic::{AtomicUsize, Ordering},
        Arc,
    },
    time::{Duration, Instant},
};

use axum::{
    extract::{Query, State},
    http::StatusCode,
    response::Response,
};
use console::style;
use indicatif::{ProgressBar, ProgressStyle};

use crate::{
    model::models::{HashRecord, SearchParams},
    AppState,
};

use super::handle_search::handle_search;

static ACTIVE_EXPORTS: AtomicUsize = AtomicUsize::new(0);

use chrono::Utc;
use serde_json;
use sha2::Digest;
use std::fs::OpenOptions;
use std::io::Write;

use std::fs::File;
use std::io::BufWriter;

fn append_hash_record(hash: String, source: String) -> std::io::Result<()> {
    let record = HashRecord {
        hash,
        timestamp: Utc::now(),
        source,
    };

    let json_line = serde_json::to_string(&record).expect("Failed to serialize record");

    let mut file = OpenOptions::new()
        .create(true)
        .append(true)
        .open("export/.hash_registry.jsonl")?;

    writeln!(file, "{}", json_line)?;
    Ok(())
}

// Handle export API request
pub async fn handle_export(
    State(state): State<Arc<AppState>>,
    Query(params): Query<SearchParams>,
) -> Result<Response, (StatusCode, String)> {
    // Check if there are too many active exports
    let active = ACTIVE_EXPORTS.load(Ordering::SeqCst);
    if active >= 1 {
        return Err((
            StatusCode::TOO_MANY_REQUESTS,
            "Too many exports in progress".to_string(),
        ));
    }

    let state_clone = state.clone();
    let params_clone = params.clone();

    tokio::spawn(async move {
        if let Err(e) = run_export(state_clone, params_clone).await {
            eprintln!("Export failed: {e}");
        }
    });

    Ok(Response::builder()
        .status(StatusCode::ACCEPTED)
        .body(axum::body::Body::from(
            "Export started. Check the export queue for progress.",
        ))
        .unwrap())
}

async fn run_export(state: Arc<AppState>, mut params: SearchParams) -> Result<(), String> {
    // Increment active exports counter
    ACTIVE_EXPORTS.fetch_add(1, Ordering::SeqCst);

    // Ensure we decrement the counter when we're done
    struct ExportGuard;
    impl Drop for ExportGuard {
        fn drop(&mut self) {
            ACTIVE_EXPORTS.fetch_sub(1, Ordering::SeqCst);
        }
    }
    let _guard = ExportGuard;

    // Get a clone of the WebSocket sender with a limited scope
    let tx = {
        let ws_state = state.ws_state.lock().await;
        ws_state.tx.clone()
    };

    // Start timing
    let start_time = Instant::now();

    // Initialize progress bar
    let pb = ProgressBar::new_spinner();
    pb.enable_steady_tick(Duration::from_millis(100));
    pb.set_style(
        ProgressStyle::default_spinner()
            .tick_chars("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
            .template("{spinner:.blue} {msg}")
            .unwrap(),
    );

    // Initial search to get total count
    pb.set_message(style("🔍 Initializing search...").cyan().to_string());
    let res = handle_search(State(state.clone()), Query(params.clone()))
        .await
        .map_err(|e| e.1)?;
    let total = res.0.total;

    if total == 0 {
        let duration = start_time.elapsed();
        pb.finish_with_message(
            style(format!("ℹ No results found (took {duration:.2?})"))
                .yellow()
                .bold()
                .to_string(),
        );
        return Ok(());
    }

    // Set up progress bar with total count
    pb.finish_with_message(
        style(format!("🔍 Found {total} logs to export..."))
            .green()
            .to_string(),
    );

    let pb = ProgressBar::new(total);
    pb.set_style(
        ProgressStyle::default_bar()
            .template("{spinner:.green} {msg}\n{wide_bar} {pos}/{len} ({eta})\n")
            .unwrap()
            .progress_chars("##-"),
    );

    // Ensure export directory exists
    let export_dir = std::path::Path::new("export");
    if !export_dir.exists() {
        std::fs::create_dir(export_dir)
            .map_err(|e| format!("Failed to create export directory: {e}"))?;
    }

    // Generate filename and save the file
    let timestamp = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .map_err(|e| format!("System time error: {e}"))?
        .as_secs();

    let ip = params.source_ip.clone().map(|ip| ip.replace('.', "_"));
    let filename = match ip {
        Some(ip) => format!("{ip}_{timestamp}.csv"),
        None => format!("all_{timestamp}.csv"),
    };
    let filepath = export_dir.join(&filename);

    let file = File::create(&filepath).map_err(|e| format!("Failed to create file: {e}"))?;
    let writer = BufWriter::new(file);

    // Create CSV writer and process data in chunks
    //let csv_data = {
    let mut wtr = csv::Writer::from_writer(writer);

    // Write CSV header
    wtr.write_record(["message"])
        .map_err(|e| format!("Failed to write CSV header: {e}"))?;

    let mut processed = 0;
    let mut last_progress_update = Instant::now();

    while processed < total {
        // Make search request
        let response = match handle_search(State(state.clone()), Query(params.clone())).await {
            Ok(r) => r,
            Err(e) => {
                let duration = start_time.elapsed();
                pb.finish_with_message(
                    style(format!("✗ Search failed after {duration:.2?}: {}", e.1))
                        .red()
                        .bold()
                        .to_string(),
                );
                return Err(format!("Search failed: {}", e.1));
            }
        };

        // Break if no more results
        if response.0.hits.is_empty() {
            break;
        }

        // Process hits
        for hit in &response.0.hits {
            //let timestamp = hit.get("timestamp").and_then(|t| t.as_str()).unwrap_or("");
            let source_ip = hit
                .get("source_ip")
                .and_then(|ip| ip.as_str())
                .unwrap_or("");
            let message = hit.get("message").and_then(|m| m.as_str()).unwrap_or("");

            if let Err(e) = wtr.write_record([message]) {
                return Err(format!("Failed to write record: {e}"));
            }

            //wtr.serialize(message);

            processed += 1;
            let elapsed = start_time.elapsed();
            let rate = processed as f64 / elapsed.as_secs_f64();
            pb.set_message(
                style(format!(
                    "📤 Exporting logs... (Last: {source_ip}, {rate:.2} logs/sec)"
                ))
                .cyan()
                .to_string(),
            );

            // Update progress periodically
            if last_progress_update.elapsed() > Duration::from_millis(100) {
                let progress = (processed as f64 / total as f64) * 100.0;
                let _ = tx.send(format!("{progress:.2}"));
                last_progress_update = Instant::now();
            }
        }

        // Update progress bar
        pb.set_position(processed);

        // Update pagination
        if let Some(last_hit) = response.0.hits.last() {
            if let Some(ts) = last_hit
                .get("index_timestamp")
                .and_then(|ts| ts.as_number())
            {
                params.end_index_timestamp = Some(ts.to_string());
            } else {
                break;
            }
        } else {
            break;
        }

        // Small delay to prevent overwhelming the server
        tokio::time::sleep(Duration::from_millis(100)).await;

        wtr.flush()
            .map_err(|e| format!("Failed to flush CSV writer: {e}"))?;
    }

    //wtr.into_inner().map_err(|e| format!("Failed to finalize CSV: {e}"))?
    //};

    // Final progress update
    let _ = tx.send("100".to_string());
    pb.finish_with_message(
        style(format!(
            "✓ Export completed in {:.2?}",
            start_time.elapsed()
        ))
        .green()
        .bold()
        .to_string(),
    );

    //std::fs::write(&filepath, csv_data)
    //    .map_err(|e| format!("Failed to write export file: {e}"))?;

    println!(
        "\n📊 {}: {} (in {:.2?})",
        style("Total records exported").bold(),
        style(pb.position()).green().bold(),
        start_time.elapsed()
    );

    println!("💾 {}", style(format!("File saved as: {filepath:?}")).dim());

    let mut reader = std::fs::File::open(&filepath)
        .map_err(|e| format!("Failed to open file {:?}: {}", filepath, e))?;

    let mut hasher = sha2::Sha256::new();

    std::io::copy(&mut reader, &mut hasher)
        .map_err(|e| format!("Failed to copy file {:?}: {}", filepath, e))?;

    let hash_string = hasher.finalize();
    let hash_string = hex::encode(hash_string);
    let hash_string = format!("SHA256:{}", hash_string);

    append_hash_record(hash_string, filename.to_string())
        .map_err(|e| format!("Failed to append hash record: {e}"))?;

    Ok(())
}
