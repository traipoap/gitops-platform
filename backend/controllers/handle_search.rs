use std::sync::Arc;

use axum::{
    extract::{Query, State},
    http::StatusCode,
    Json,
};

use crate::{
    model::models::{SearchParams, SearchResponse},
    AppState,
};

use dotenvy::dotenv;
use std::env;

// Handle search API request
pub async fn handle_search(
    State(state): State<Arc<AppState>>,
    Query(params): Query<SearchParams>,
) -> Result<Json<SearchResponse>, (StatusCode, String)> {
    dotenv().ok();
    let quickwit_url = env::var("QUICKWIT_URL").expect("QUICKWIT_URL must be set");

    // Build the Quickwit query
    let mut query_parts = Vec::new();

    if let Some(sts) = params.from_timestamp {
        if let Some(ets) = params.to_timestamp {
            query_parts.push(format!("timestamp:[{sts} TO {ets}]"));
        }
    }

    if let Some(ip) = params.source_ip {
        query_parts.push(format!("source_ip:{ip}"));
    }
    if let Some(filter) = params.custom_filter {
        query_parts.push(format!("message:{filter}"));
    }

    if let Some(end_index_timestamp) = params.end_index_timestamp {
        query_parts.push(format!("index_timestamp:[* TO {}]", end_index_timestamp));
    } else {
        query_parts.push("index_timestamp:[* TO *]".to_string());
    }

    let query = if query_parts.is_empty() {
        "*".to_string()
    } else {
        query_parts.join(" AND ")
    };

    // Call Quickwit API
    let response = state
        .client
        .post(quickwit_url + "/api/v1/syslogs/search")
        .json(&serde_json::json!({
            "query": query,
            "max_hits": params.limit.unwrap_or(10000),
            "sort_by_field": "index_timestamp",
        }))
        .send()
        .await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;

    if !response.status().is_success() {
        return Err((
            StatusCode::from_u16(response.status().as_u16())
                .unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            format!("Quickwit API error: {}", response.status()),
        ));
    }

    let result: serde_json::Value = response
        .json()
        .await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;

    let hits = result["hits"].as_array().cloned().unwrap_or_default();

    let total = result["num_hits"].as_u64().unwrap_or(hits.len() as u64);

    Ok(Json(SearchResponse { hits, total }))
}
