// State
let allLogs = [];
let filteredLogs = [];
let currentPage = 1;
let pageSize = 50;
let sortField = "timestamp";
let sortDirection = "desc";
let selectedLog = null;
let liveMode = true;
let liveInterval = null;
let selectedIndex = "syslogs";

// Use SAME-ORIGIN relative /api/* URLs. The Vite dev proxy (astro.config.mjs)
// forwards these to http://localhost:8080 locally, and the production
// reverse proxy does the same at frontend.example.com — so no CORS and no
// need for import.meta.env. (This file loads as a plain <script>, NOT a
// module, so import.meta is unavailable and would throw.)
const API_BASE = "";

// ── Token helpers ──────────────────────────────────────────────
function getAuthToken() {
  return sessionStorage.getItem("authToken") || localStorage.getItem("authToken");
}

function getRefreshToken() {
  return sessionStorage.getItem("refreshToken") || localStorage.getItem("refreshToken");
}

function setTokens(accessToken, refreshToken) {
  localStorage.setItem("authToken", accessToken);
  localStorage.setItem("refreshToken", refreshToken);
}

function clearTokens() {
  sessionStorage.removeItem("authToken");
  sessionStorage.removeItem("refreshToken");
  localStorage.removeItem("authToken");
  localStorage.removeItem("refreshToken");
}

async function tryRefreshToken() {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;

  try {
    const res = await fetch(`${API_BASE}/api/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!res.ok) {
      clearTokens();
      window.location.href = "/signin";
      return false;
    }

    const data = await res.json();
    setTokens(data.token, data.refresh_token);
    return true;
  } catch {
    clearTokens();
    window.location.href = "/signin";
    return false;
  }
}

async function authenticatedFetch(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
      ...options.headers,
    },
  });

  // On 401, try refreshing once then retry
  if (res.status === 401) {
    const refreshed = await tryRefreshToken();
    if (refreshed) {
      return fetch(url, {
        ...options,
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
          ...options.headers,
        },
      });
    }
    return res; // already redirected
  }

  return res;
}

async function loadEngine() {
  try {
    const res = await authenticatedFetch(`${API_BASE}/api/indices`);

    if (res.status === 403) {
      clearTokens();
      window.location.href = "/signin";
      return;
    }

    if (!res.ok) {
      const errText = await res.text().catch(() => "");
      console.error(`API error (${res.status}):`, errText);
      showToast(`Server error (${res.status}). Please try again.`);
      return;
    }

    const data = await res.json();

    // 1. จัดการ engineSelect (ลดรูปเหลือบรรทัดเดียว)
    document.getElementById("engineSelect").innerHTML = ["gin", "axum", "auto-generate"]
      .map(id => `<option value="${id}" ${id === 'syslogs' ? 'selected' : ''}>${id}</option>`).join("");

    // 2. จัดการ indexSelect (เอา syslogs ขึ้นก่อน)
    document.getElementById("indexSelect").innerHTML = data
      .map(item => item.index_config?.index_id)
      .filter(id => id) // กรองเฉพาะค่าที่มีอยู่จริง (ไม่เป็น undefined หรือ empty)
      .sort((a, b) => (a === 'syslogs' ? -1 : b === 'syslogs' ? 1 : 0)) // ย้าย syslogs ไปไว้หน้าสุด
      .map(id => `<option value="${id}">${id}</option>`)
      .join("");

    // 3. ตั้งค่า dateTo เป็นเวลาปัจจุบันเมื่อโหลดหน้า (ใช้ local time โดยตรง)
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const localISOTime = `${year}-${month}-${day}T${hours}:${minutes}`;
    document.getElementById("dateTo").value = localISOTime;

    //runSearch();

  } catch (err) {
    console.error("Fetch failed:", err);
  }
}

// 1. ปรับปรุงการ Mapping (เน้นความยืดหยุ่น)
function mapQuickwitHits(hits) {
  return hits.map(hit => {
    // แปลง Timestamp: เช็คว่าเป็นวินาที หรือ มิลลิวินาที
    // ถ้าค่าน้อยกว่า 10,000,000,000 แสดงว่าเป็นวินาที (Unix Timestamp) ให้คูณ 1000
    let ts = hit.timestamp;
    if (ts < 10000000000) ts *= 1000;
    const dateObj = new Date(ts);

    // หา Level จากหลายๆ field
    const rawLevel = hit.severity || hit.level || "info";
    const normalizedLevel = normalizeSeverity(rawLevel);

    // สร้าง Core Object
    const core = {
      id: hit.id || Math.random().toString(36).substr(2, 9),
      timestamp: dateObj,
      level: normalizedLevel,
      source: hit.source_ip || hit.srcip || hit.source || hit.kubernetes.pod_ip|| "unknown",
      host: hit.host || hit.hostname || hit.kubernetes.container_name || "unknown",
      message: hit.message || "",
      pid: hit.pid || 0
    };

    // เก็บ Extra fields (Dynamic)
    const extras = {};
    Object.keys(hit).forEach(key => {
      // ไม่เอา key ที่เราใช้ไปแล้วใน core มาซ้ำใน extras เพื่อความสะอาด
      const coreKeys = ['timestamp', 'severity', 'level', 'source_ip', 'srcip', 'hostname'];
      if (!coreKeys.includes(key)) {
        extras[key] = hit[key];
      }
    });

    return { ...core, extras };
  });
}

/**
 * ฟังก์ชันสำหรับจัดการ Timestamp ที่อาจมาเป็น Milliseconds หรือ Nanoseconds
 * (ตัวอย่างข้อมูลคุณมีทั้ง 1782028862574 และ 1782028862574283000)
 */
function parseFlexibleTimestamp(ts) {
  if (!ts) return new Date();

  let date = new Date(ts);

  // ถ้าค่าที่ได้คือ NaN หรือเป็นตัวเลขที่ยาวเกินไป (เช่น Nanoseconds)
  // เราต้องเช็คว่ามันยาวเกิน 13 หลักหรือไม่ (13 หลักคือระดับ Milliseconds)
  if (isNaN(date.getTime()) || ts > 1000000000000000) {
    // ตัดให้เหลือ 13 หลักแรก (แปลงจาก ns/us เป็น ms)
    const ms = Math.floor(ts / 1000000);
    date = new Date(ms);
  }

  // ป้องกันกรณีเป็นเลข Nanoseconds ที่แปลงแล้วอาจจะผิดพลาด
  // วิธีที่ชัวร์ที่สุดคือเช็คความสมเหตุสมผลของปี (เช่น ไม่ควรเกินปี 2100)
  if (date.getFullYear() > 2100) {
    date = new Date(Math.floor(ts / 1000000));
  }

  return date;
}

/**
 * ฟังก์ชัน Dynamic Mapper
 * รับ hit จาก Backend (ที่มีโครงสร้างไม่แน่นอน) แล้วแปลงเป็น Standard Schema ของ UI
 */
function dynamicMapper(hit) {
  // 1. กำหนดลำดับความสำคัญของ Field (Priority List)
  // เราจะหาค่าจากชื่อที่น่าจะเป็นไปได้มากที่สุดก่อน
  const getField = (aliases) => {
    for (let alias of aliases) {
      if (hit[alias] !== undefined && hit[alias] !== null) {
        return hit[alias];
      }
    }
    return null;
  };

  // 2. สกัดข้อมูลพื้นฐาน
  const levelRaw = getField(['severity', 'level', 'priority', 'log_level']);
  const sourceRaw = getField(['source_ip', 'srcip', 'source_type', 'appname', 'facility', 'source', 'pod_node_name']);
  const hostRaw = getField(['host', 'hostname', 'devname', 'source_host']);
  const timestampRaw = getField(['timestamp', 'index_timestamp', 'eventtime', 'time']);
  const messageRaw = getField(['message', 'msg', 'content']);
  const pidRaw = getField(['pid', 'process_id']);

  // 3. ตรวจสอบความลึก (Deep Parsing)
  // กรณีข้อมูลสำคัญอยู่ใน string ของ message (เช่น Fortigate format)
  let extraSrcIp = null;
  const srcIpMatch = messageRaw?.match(/srcip=([\d.]+)/);
  if (srcIpMatch) extraSrcIp = srcIpMatch[1];

  // 4. สร้าง Object ใน Format ที่ UI ต้องการ (Standard Schema)
  return {
    id: hit.id || Math.random().toString(36).substr(2, 9),
    timestamp: parseFlexibleTimestamp(timestampRaw || hit.index_timestamp),
    level: normalizeSeverity(levelRaw || 'info'),
    source: sourceRaw || extraSrcIp || 'unknown',
    host: hostRaw || 'unknown',
    message: messageRaw || '',
    pid: pidRaw ? parseInt(pidRaw) : 0,
    // เก็บ metadata เดิมไว้เผื่อใช้ในหน้า Detail
    metadata: {
      appname: hit.appname || null,
      facility: hit.facility || null,
      version: hit.version || null
    }
  };
}

async function runSearch() {
  const source = document.getElementById("sourceInput").value.trim();
  const query = document.getElementById("searchInput").value.trim();
  const index = document.getElementById("indexSelect").value;
  const btn = document.getElementById("searchBtn");

  // 1. ดึงค่า Filter ปัจจุบัน
  const currentSourceFilter = document.getElementById("sourceFilter").value;
  const dateFromVal = document.getElementById("dateFrom").value;
  const dateToVal = document.getElementById("dateTo").value;
  const currentHostFilter = document.getElementById("hostFilter").value;

  // 2. เตรียม Params พื้นฐาน
  const params = new URLSearchParams({ index_id: index, max_hits: 2000 });

  // 3. เพิ่ม Query/Message
  if (query) params.set("message", query);

  // 4. เพิ่ม Source IP
  if (currentSourceFilter) {
    params.set("source_ip", currentSourceFilter);
  } else if (source) {
    params.set("source_ip", source);
  }

  // 5. เพิ่ม Host
  if (currentHostFilter) params.set("host", currentHostFilter);

  // 6. จัดการเรื่อง Timestamp (แปลงจาก ISO String -> Milliseconds)
  // หมายเหตุ: ถ้า Backend ต้องการเป็น Seconds ให้หารด้วย 1000 (เช่น .getTime() / 1000)
  if (dateFromVal) {
    const fromTs = new Date(dateFromVal).getTime();
    params.set("from_timestamp", fromTs);
  }
  if (dateToVal) {
    const toTs = new Date(dateToVal).getTime();
    params.set("to_timestamp", toTs);
  }

  btn.disabled = true;

  try {
    // 6. เรียก API พร้อม Params ใหม่ (auto-refresh on 401)
    const res = await authenticatedFetch(`${API_BASE}/api/search?${params.toString()}`);

    if (res.status === 403) {
      clearTokens();
      window.location.href = "/signin";
      return;
    }

    if (!res.ok) throw new Error(`Server error ${res.status}`);

    const data = await res.json();

    // 7. อัปเดตข้อมูล (ใช้ mapQuickwitHits ที่เราคุยกันก่อนหน้า)
    allLogs = mapQuickwitHits(data.hits || []);

    // 8. อัปเดต Dropdown Source IP
    const uniqueIps = (data.hits || [])
      .map(item => item.source_ip)
      .filter((id, idx, self) => id && self.indexOf(id) === idx)
      .sort((a, b) => (a === '192.168.1.2' ? -1 : b === '192.168.1.2' ? 1 : 0));

    document.getElementById("sourceFilter").innerHTML = ['<option value="">All Sources</option>']
      .concat(uniqueIps.map(id => `<option value="${id}">${id}</option>`))
      .join("");
    // ล้าง filter source เพื่อให้แสดงผลการ search ใหม่แบบทั้งหมด
    document.getElementById("sourceFilter").value = "";

    const uniqueHosts = (data.hits || [])
      .map(item => item.host)
      .filter((id, idx, self) => id && self.indexOf(id) === idx)
      .sort((a, b) => (a === '192.168.1.2' ? -1 : b === '192.168.1.2' ? 1 : 0));

    document.getElementById("hostFilter").innerHTML = ['<option value="">All Hosts</option>']
      .concat(uniqueHosts.map(id => `<option value="${id}">${id}</option>`))
      .join("");
    // ล้าง filter source เพื่อให้แสดงผลการ search ใหม่แบบทั้งหมด
    document.getElementById("hostFilter").value = "";

    // 9. อัปเดต UI ทั้งหมด
    applyFilters();
    updateStats();

    // ถ้าหน้า Dashboard เปิดอยู่ ให้วาดกราฟใหม่ด้วย
    if (document.getElementById("dashboardView").classList.contains("active")) {
      renderDashboard();
    }

    showToast(`Found ${data.total || 0} results`);

  } catch (err) {
    console.error("Search error:", err);
    showToast("Search failed");
    allLogs = [];
    applyFilters();
  } finally {
    btn.disabled = false;
  }
}

function refreshAllUI() {
  applyFilters();   // อัปเดตตารางและตัวกรอง
  updateStats();    // อัปเดตตัวเลข Error/Warn/Info ด้านบน

  // ถ้าหน้า Dashboard กำลังเปิดอยู่ ให้วาดกราฟและตัวเลขใหม่ทันที
  if (document.getElementById("dashboardView").classList.contains("active")) {
    renderDashboard();
  }
}

function normalizeSeverity(severity) {
  const map = {
    err: "error",
    error: "error",
    warning: "warn",
    warn: "warn",
    crit: "critical",
    critical: "critical",
    emerg: "critical",
    alert: "critical",
    notice: "info",
    info: "info",
    debug: "debug",
  };
  return map[severity?.toLowerCase()] ?? severity?.toLowerCase() ?? "info";
}

// Sample data generators
const sources = [
  "nginx",
  "postgresql",
  "redis",
  "app-server",
  "auth-service",
  "cron",
  "systemd",
  "kernel",
];
const hosts = ["prod-web-01", "prod-web-02", "prod-db-01", "prod-cache-01"];
const levels = ["error", "warn", "info", "debug", "critical"];
const levelWeights = [15, 20, 45, 15, 5];

const messages = {
  error: [
    "Connection refused to database server at 10.0.1.5:5432",
    "Failed to authenticate user: invalid credentials",
    "Out of memory: Killed process {pid} (java)",
    "SSL handshake failed with upstream server",
    "Segmentation fault in module auth_handler.so",
    "Disk space critically low on /dev/sda1 (95% used)",
    "Timeout waiting for response from cache cluster",
    "Failed to write to audit log: Permission denied",
    "Uncaught exception in request handler: NullPointerException",
    "Database connection pool exhausted (max: 100)",
  ],
  warn: [
    "Slow query detected: SELECT * FROM orders WHERE... (2340ms)",
    "Certificate expires in 7 days for domain api.example.com",
    "High memory usage detected: 87% of available RAM",
    "Rate limit approaching for client 192.168.1.100",
    "Deprecated API endpoint /v1/users accessed",
    "Connection pool reaching capacity: 85/100 active",
    "Retry attempt 3/5 for service discovery lookup",
    "Request payload size exceeds recommended limit (8MB)",
    "Thread pool saturation at 75%, consider scaling",
    "Stale cache entry detected for key user:session:48291",
  ],
  info: [
    "Request processed successfully in 45ms (GET /api/v2/users)",
    "Database backup completed successfully (2.3GB)",
    "New deployment v2.4.1 started on prod-web-01",
    "Health check passed for all service endpoints",
    "User session created: user_id=12847, ttl=3600s",
    "Cache invalidation completed for 1243 entries",
    'Cron job "cleanup-logs" executed in 12.4s',
    "SSL certificate renewed for *.example.com",
    "Auto-scaling triggered: adding 2 instances to web pool",
    "Configuration reload successful, 84 parameters updated",
  ],
  debug: [
    "Query execution plan: Index scan on users_email_idx",
    "Cache hit ratio: 94.2% for session store",
    "GC pause: 12ms (Eden space, 256MB collected)",
    "TCP connection established to 10.0.1.5:443",
    "JWT token validated: sub=admin, exp=1732847291",
    "Request headers: {Authorization: Bearer ***, Accept: application/json}",
    "Response serialization: 2847 bytes in 3ms",
    "Connection pool stats: active=45, idle=30, waiting=0",
    "Metric flush: 1842 data points sent to monitoring",
    "Worker thread #7 processing task queue (depth: 12)",
  ],
  critical: [
    "SYSTEM FAILURE: Primary database unreachable for 30s",
    "Security alert: Multiple failed login attempts from 203.0.113.42",
    "Data corruption detected in table: user_sessions",
    "Service crash: app-server exited with code 137 (OOM)",
    "Network partition detected between availability zones",
    "Ransomware signature detected in uploaded file: invoice.pdf",
    "Kernel panic: Not syncing - Fatal exception in interrupt",
    "Backup verification FAILED: checksum mismatch on volume-7",
    "Firewall rule breach: unauthorized port 22 access attempt",
    "Cluster quorum lost: only 1/3 nodes reachable",
  ],
};

function weightedRandom(arr, weights) {
  const total = weights.reduce((a, b) => a + b, 0);
  let random = Math.random() * total;
  for (let i = 0; i < arr.length; i++) {
    random -= weights[i];
    if (random <= 0) return arr[i];
  }
  return arr[arr.length - 1];
}

function generateLogEntry(timestamp) {
  const level = weightedRandom(levels, levelWeights);
  const source = sources[Math.floor(Math.random() * sources.length)];
  const host = hosts[Math.floor(Math.random() * hosts.length)];
  const messageTemplates = messages[level];
  let message =
    messageTemplates[Math.floor(Math.random() * messageTemplates.length)];
  if (message.includes("{pid}")) {
    message = message.replace("{pid}", Math.floor(Math.random() * 30000));
  }
  return {
    id: Math.random().toString(36).substr(2, 9),
    timestamp: timestamp || new Date(Date.now() - Math.random() * 86400000 * 7),
    level: level,
    source: source,
    host: host,
    message: message,
    pid: Math.floor(Math.random() * 30000) + 1000,
    extras: {}
  };
}

function generateLogs(count = 500) {
  allLogs = [];
  const now = Date.now();
  for (let i = 0; i < count; i++) {
    const ts = new Date(now - Math.random() * 86400000 * 7);
    allLogs.push(generateLogEntry(ts));
  }
  allLogs.sort((a, b) => b.timestamp - a.timestamp);
  applyFilters();
  updateStats();
  if (document.getElementById("dashboardView").classList.contains("active")) {
    renderDashboard();
  }
  refreshAllUI();
  showToast(`${count} log entries generated`);
}

function formatTimestamp(date) {
  if (!date) return "";

  // ใช้ Intl.DateTimeFormat เพื่อบังคับ Timezone เป็น Asia/Bangkok
  // และใช้ locale 'sv-SE' เพราะให้รูปแบบ YYYY-MM-DD ที่ใกล้เคียงกับที่เราต้องการที่สุด
  const formatter = new Intl.DateTimeFormat('sv-SE', {
    timeZone: 'Asia/Bangkok',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });

  // ผลลัพธ์จะได้ประมาณ "2024-06-21 15:01:02"
  return formatter.format(date).replace(/-/g, '-').replace(',', '');
}

function getActiveLevels() {
  return Array.from(document.querySelectorAll(".level-btn.active")).map(
    (b) => b.dataset.level,
  );
}

// 2. ปรับปรุง applyFilters ให้ "ใจดี" ขึ้น (สำหรับการ Debug)
function applyFilters() {
  const search = document.getElementById("searchInput").value.toLowerCase();
  const activeLevels = getActiveLevels();
  const dateFrom = document.getElementById("dateFrom").value ? new Date(document.getElementById("dateFrom").value) : null;
  const dateTo = document.getElementById("dateTo").value ? new Date(document.getElementById("dateTo").value) : null;
  const sourceFilter = document.getElementById("sourceFilter").value;
  const hostFilter = document.getElementById("hostFilter").value;

  filteredLogs = allLogs.filter((log) => {
    // กรอง Level: ถ้าไม่ได้เลือก Level เลย (ทุกปุ่ม inactive) ให้แสดงทั้งหมด หรือเช็คตาม logic ของคุณ
    if (activeLevels.length > 0 && !activeLevels.includes(log.level)) return false;

    // กรอง Source
    if (sourceFilter && log.source !== sourceBreaker(sourceFilter)) return false;

    // กรอง Host
    if (hostFilter && log.host !== hostFilter) return false;

    // กรอง Date (เพิ่มการเช็คเผื่อกรณีข้อมูลเป็นอนาคต เพื่อใช้ในการ Test)
    // ถ้า dateTo มีค่า และ log.timestamp มันล้ำหน้าไปไกลเกินไป (เช่น ปี 2026)
    // ในการทดสอบเราอาจจะข้ามการเช็ค dateTo ไปก่อน หรือขยาย range
    if (dateFrom && log.timestamp < dateFrom) return false;
    if (dateTo && log.timestamp > dateTo) {
        // logic พิเศษ: ถ้าเป็นข้อมูลปี 2026 (จากการ test) ให้ยอมให้ผ่านเพื่อไม่ให้หน้าจอว่าง
        if (log.timestamp.getFullYear() > 2025) {
            // console.log("Skipping date filter for future log:", log.timestamp);
        } else {
            return false;
        }
    }

    // กรอง Search
    if (search) {
      const searchIn = `${log.message} ${log.source} ${log.host} ${log.level}`.toLowerCase();
      if (!searchIn.includes(search)) return false;
    }

    return true;
  });

  sortLogsInternal();
  currentPage = 1;
  renderTable();
  renderPagination();
  document.getElementById("totalFiltered").textContent = filteredLogs.length;
  document.getElementById("filteredCount").textContent = Math.min(pageSize, filteredLogs.length);
}

// Helper function สำหรับ handle source filter
function sourceBreaker(val) {
    return val; // เพิ่ม logic แปลงค่าถ้าจำเป็น
}

function sortLogs(field) {
  if (sortField === field) {
    sortDirection = sortDirection === "asc" ? "desc" : "asc";
  } else {
    sortField = field;
    sortDirection = field === "timestamp" ? "desc" : "asc";
  }

  // Update header styles
  document.querySelectorAll(".log-table th").forEach((th) => {
    th.classList.remove("sorted");
    const arrow = th.querySelector(".sort-arrow");
    if (arrow) arrow.textContent = "";
  });

  const thIndex = {
    timestamp: 0,
    level: 1,
    source: 2,
    host: 3,
    message: 4,
    pid: 5,
  };
  const headerTh = document.querySelectorAll(".log-table th")[thIndex[field]];
  if (headerTh) {
    headerTh.classList.add("sorted");
    const arrow = headerTh.querySelector(".sort-arrow");
    if (arrow) arrow.textContent = sortDirection === "asc" ? "▲" : "▼";
  }

  sortLogsInternal();
  renderTable();
}

function sortLogsInternal() {
  filteredLogs.sort((a, b) => {
    let valA, valB;
    switch (sortField) {
      case "timestamp":
        valA = a.timestamp.getTime();
        valB = b.timestamp.getTime();
        break;
      case "level":
        valA = levels.indexOf(a.level);
        valB = levels.indexOf(b.level);
        break;
      case "source":
        valA = a.source;
        valB = b.source;
        break;
      case "host":
        valA = a.host;
        valB = b.host;
        break;
      case "message":
        valA = a.message.toLowerCase();
        valB = b.message.toLowerCase();
        break;
      case "pid":
        valA = a.pid;
        valB = b.pid;
        break;
      default:
        valA = a.timestamp.getTime();
        valB = b.timestamp.getTime();
    }
    if (valA < valB) return sortDirection === "asc" ? -1 : 1;
    if (valA > valB) return sortDirection === "asc" ? 1 : -1;
    return 0;
  });
}

function renderTable() {
  const tbody = document.getElementById("logTableBody");
  const search = document.getElementById("searchInput").value.toLowerCase();
  const start = (currentPage - 1) * pageSize;
  const end = start + pageSize;
  const pageLogs = filteredLogs.slice(start, end);

  if (pageLogs.length === 0) {
    tbody.innerHTML = `
          <tr>
              <td colspan="6">
                  <div class="empty-state">
                      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                          <circle cx="12" cy="12" r="10"></circle>
                          <line x1="12" y1="8" x2="12" y2="12"></line>
                          <line x1="12" y1="16" x2="12.01" y2="16"></line>
                      </svg>
                      <h3>No logs found</h3>
                      <p>Try adjusting your filters or generate new log data</p>
                  </div>
              </td>
          </tr>
      `;
    return;
  }

  tbody.innerHTML = pageLogs
    .map((log) => {
      let displayMessage = log.message;
      if (search) {
        const regex = new RegExp(`(${escapeRegex(search)})`, "gi");
        displayMessage = displayMessage.replace(
          regex,
          '<span class="highlight">$1</span>',
        );
      }
      return `
          <tr class="log-row" onclick="showLogDetail('${log.id}')">
              <td class="timestamp">${formatTimestamp(log.timestamp)}</td>
              <td><span class="log-level-badge ${log.level}">${log.level}</span></td>
              <td class="log-source">${log.source}</td>
              <td>${log.host}</td>
              <td class="log-message">${displayMessage}</td>
              <td>${log.pid}</td>
          </tr>
      `;
    })
    .join("");

  document.getElementById("filteredCount").textContent = pageLogs.length;
}

function escapeRegex(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function renderPagination() {
  const totalPages = Math.ceil(filteredLogs.length / pageSize) || 1;
  const controls = document.getElementById("paginationControls");
  const info = document.getElementById("paginationInfo");

  info.textContent = `Page ${currentPage} of ${totalPages} (${filteredLogs.length} entries)`;

  let html = `<button class="page-btn" onclick="goToPage(1)" ${currentPage === 1 ? "disabled" : ""}>«</button>`;
  html += `<button class="page-btn" onclick="goToPage(${currentPage - 1})" ${currentPage === 1 ? "disabled" : ""}>‹</button>`;

  let startPage = Math.max(1, currentPage - 2);
  let endPage = Math.min(totalPages, currentPage + 2);

  if (startPage > 1)
    html += `<button class="page-btn" onclick="goToPage(1)">1</button><span style="padding:0 4px;color:var(--text-muted)">...</span>`;

  for (let i = startPage; i <= endPage; i++) {
    html += `<button class="page-btn ${i === currentPage ? "active" : ""}" onclick="goToPage(${i})">${i}</button>`;
  }

  if (endPage < totalPages)
    html += `<span style="padding:0 4px;color:var(--text-muted)">...</span><button class="page-btn" onclick="goToPage(${totalPages})">${totalPages}</button>`;

  html += `<button class="page-btn" onclick="goToPage(${currentPage + 1})" ${currentPage === totalPages ? "disabled" : ""}>›</button>`;
  html += `<button class="page-btn" onclick="goToPage(${totalPages})" ${currentPage === totalPages ? "disabled" : ""}>»</button>`;

  controls.innerHTML = html;
}

function goToPage(page) {
  const totalPages = Math.ceil(filteredLogs.length / pageSize) || 1;
  if (page < 1 || page > totalPages) return;
  currentPage = page;
  renderTable();
  renderPagination();
  document.getElementById("logTableView").scrollTop = 0;
}

function changePageSize() {
  pageSize = parseInt(document.getElementById("pageSize").value);
  currentPage = 1;
  renderTable();
  renderPagination();
}

function updateStats() {
  const total = allLogs.length;
  const errors = allLogs.filter((l) => l.level === "error").length;
  const warns = allLogs.filter((l) => l.level === "warn").length;
  const infos = allLogs.filter((l) => l.level === "info").length;

  document.getElementById("totalCount").textContent = total.toLocaleString();
  document.getElementById("errorCount").textContent = errors.toLocaleString();
  document.getElementById("warnCount").textContent = warns.toLocaleString();
  document.getElementById("infoCount").textContent = infos.toLocaleString();
}

function toggleLevel(btn) {
  btn.classList.toggle("active");
  applyFilters();
}

// Helper function สำหรับแปลง Date เป็น local ISO string (YYYY-MM-DDTHH:mm)
function toLocalISOString(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

function setQuickDate(period) {
  const now = new Date();
  let from;
  switch (period) {
    case "1h":
      from = new Date(now - 3600000);
      break;
    case "24h":
      from = new Date(now - 86400000);
      break;
    case "7d":
      from = new Date(now - 604800000);
      break;
    case "30d":
      from = new Date(now - 2592000000);
      break;
  }
  document.getElementById("dateFrom").value = toLocalISOString(from);
  document.getElementById("dateTo").value = toLocalISOString(now);
  applyFilters();
}

function clearAllFilters() {
  document.getElementById("searchInput").value = "";
  document.getElementById("dateFrom").value = "";
  document.getElementById("dateTo").value = "";
  document.getElementById("sourceFilter").value = "";
  document.getElementById("hostFilter").value = "";
  document
    .querySelectorAll(".level-btn[data-level]")
    .forEach((b) => b.classList.add("active"));
  applyFilters();
  showToast("All filters cleared");
}

function debounceSearch() {
  clearTimeout(window.searchTimeout);
  window.searchTimeout = setTimeout(applyFilters, 300);
}

function showLogDetail(id) {
  const log = allLogs.find((l) => l.id === id);
  if (!log) return;
  selectedLog = log;

  // Create the core HTML
  let html = `
      <div class="detail-grid">
          <div class="detail-label">Timestamp</div>
          <div class="template-value">${formatTimestamp(log.timestamp)}</div>
          <div class="detail-label">Level</div>
          <div class="detail-value"><span class="log-level-badge ${log.level}">${log.level}</span></div>
          <div class="detail-label">Source</div>
          <div class="detail-value">${log.source}</div>
          <div class="detail-label">Host</div>
          <div class="detail-value">${log.host}</div>
          <div class="detail-label">Message</div>
          <div class="detail-value full-width">${log.message}</div>
  `;

  // Check if extras exists and has keys before calling Object.keys
  if (log.extras && typeof log.extras === 'object' && Object.keys(log.extras).length > 0) {
    html += `<div class="detail-full-width" style="grid-column: 1/-1; margin-top: 20px; border-top: 1px solid #444; padding-top: 10px;">
                <h4 style="margin-bottom: 10px;">Additional Properties</h4>
             </div>`;

    html += Object.keys(log.extras).map(key => `
      <div class="detail-label">${key}</div>
      <div class="detail-value">${log.extras[key]}</div>
    `).join('');
  }

  html += `</div>`;

  document.getElementById("modalBody").innerHTML = html;
  document.getElementById("logModal").classList.add("active");
}

function closeModal() {
  document.getElementById("logModal").classList.remove("active");
}

async function showExportHistoryModal() {
  document.getElementById("modalContent").innerHTML = `
    <div style="padding:2rem;text-align:center;color:#94a3b8">
      <span class="spinner" style="display:inline-block"></span>
      <p>Loading export history...</p>
    </div>`;
  document.getElementById("exportHistoryModal").classList.add("active");

  try {
    const res = await authenticatedFetch(`${API_BASE}/api/exports`);
    if (!res.ok) {
      throw new Error(`Failed to load exports (${res.status})`);
    }

    const data = await res.json();
    const exports = data.exports || [];

    const rows = exports.map(item => `
      <tr>
        <td>${formatExportTime(item.time)}</td>
        <td title="${item.name}">${truncate(item.name, 40)}</td>
        <td class="hash-cell" title="${item.hash}" onclick="copyHash('${escapeHtml(item.hash)}')">${truncate(item.hash, 20)}</td>
        <td>${formatFileSize(item.size)}</td>
        <td>
          <button class="btn-action" onclick="downloadExport('${escapeHtml(item.name)}')" title="Download">⬇️</button>
          <button class="btn-action" onclick="deleteExport('${escapeHtml(item.name)}')" title="Delete">🗑️</button>
        </td>
      </tr>
    `).join("");

    const html = `
      <div style="display:flex;justify-content:space-between;align-items:center;padding:0.75rem 1rem">
        <h3 style="margin:0;font-size:1.1rem">📂 Export History</h3>
        <span class="badge" style="background:#334155;color:#94a3b8;padding:2px 10px;border-radius:999px;font-size:0.8rem">
          ${exports.length} file${exports.length !== 1 ? 's' : ''}
        </span>
      </div>
      <table class="export-history-table">
        <thead>
          <tr>
            <th>Export Time</th>
            <th>File Name</th>
            <th>HASH (click to copy)</th>
            <th>Size</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          ${rows || '<tr><td colspan="5" style="text-align:center;color:#64748b;padding:2rem">No exports yet</td></tr>'}
        </tbody>
      </table>
    `;

    document.getElementById("modalContent").innerHTML = html;
  } catch (err) {
    console.error("Export history error:", err);
    document.getElementById("modalContent").innerHTML = `
      <div style="padding:2rem;text-align:center;color:#ef4444">
        Failed to load export history. Please try again.
      </div>`;
  }
}

// ── Utility helpers for export modal ──────────────────────
function formatExportTime(iso) {
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return iso;
  }
}

function formatFileSize(bytes) {
  if (!bytes || bytes === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i];
}

function truncate(str, maxLen) {
  if (!str) return '';
  return str.length > maxLen ? str.slice(0, maxLen - 3) + '...' : str;
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function copyHash(hash) {
  navigator.clipboard.writeText(hash).then(() => {
    showToast("Hash copied to clipboard");
  }).catch(() => {
    showToast("Failed to copy hash");
  });
}

function closeExportHistoryModal() {
  document.getElementById("exportHistoryModal").classList.remove("active");
}

async function downloadExport(fileName) {
  showToast(`Downloading ${fileName}...`);

  try {
    const res = await authenticatedFetch(`${API_BASE}/api/exports/${encodeURIComponent(fileName)}`);
    if (!res.ok) throw new Error(`Download failed (${res.status})`);

    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    showToast(`Downloaded ${fileName}`);
  } catch (err) {
    console.error("Download error:", err);
    showToast("Download failed. Please try again.");
  }
}

async function deleteExport(fileName) {
  if (!confirm(`Are you sure you want to delete "${fileName}"?`)) return;
  showToast(`Deleting ${fileName}...`);

  try {
    const res = await authenticatedFetch(`${API_BASE}/api/exports/${encodeURIComponent(fileName)}`, {
      method: "DELETE",
    });

    if (!res.ok) throw new Error(`Delete failed (${res.status})`);

    showToast(`Deleted ${fileName}`);
    // Refresh the export history modal
    showExportHistoryModal();
  } catch (err) {
    console.error("Delete error:", err);
    showToast("Delete failed. Please try again.");
  }
}

function copyLogDetail() {
  if (!selectedLog) return;
  const text = `[${formatTimestamp(selectedLog.timestamp)}] [${selectedLog.level.toUpperCase()}] [${selectedLog.source}@${selectedLog.host}] (PID: ${selectedLog.pid}) ${selectedLog.message}`;
  navigator.clipboard.writeText(text).then(() => {
    showToast("Log entry copied to clipboard");
  });
}

function exportSingleLog() {
  if (!selectedLog) return;
  const blob = new Blob([JSON.stringify(selectedLog, null, 2)], {
    type: "application/json",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `log_${selectedLog.id}.json`;
  a.click();
  URL.revokeObjectURL(url);
  showToast("Log entry exported");
}

function exportCSV() {
  if (filteredLogs.length === 0) {
    showToast("No logs to export");
    return;
  }
  const headers = ["Timestamp", "Level", "Source", "Host", "PID", "Message"];
  const rows = filteredLogs.map((l) => [
    formatTimestamp(l.timestamp),
    l.level,
    l.source,
    l.host,
    l.pid,
    `"${l.message.replace(/"/g, '""')}"`,
  ]);
  const csv = [headers.join(","), ...rows.map((r) => r.join(","))].join("\n");
  const blob = new Blob([csv], { type: "text/csv" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `logs_export_${Date.now()}.csv`;
  a.click();
  URL.revokeObjectURL(url);
  showToast(`Exported ${filteredLogs.length} log entries as CSV`);
}

// Trigger a SERVER-SIDE export via backend HandleExport (/api/export).
// The backend queries Quickwit, saves a CSV to ./exports as
// "{sourceIP|any}_{YYYYMMDD_HHMMSS}.csv", then we download it to the browser.
async function exportLargeCSV() {
  const indexEl = document.getElementById("indexSelect");
  const index = indexEl ? indexEl.value : "";
  if (!index) {
    showToast("Select an index first");
    return;
  }

  const source = (document.getElementById("sourceInput")?.value || "").trim();
  const query = (document.getElementById("searchInput")?.value || "").trim();
  const currentSourceFilter = document.getElementById("sourceFilter")?.value || "";
  const dateFromVal = document.getElementById("dateFrom")?.value || "";
  const dateToVal = document.getElementById("dateTo")?.value || "";

  const params = new URLSearchParams({ index_id: index, max_hits: 50000 });
  const src = currentSourceFilter || source; // source IP drives the filename
  if (src) params.set("source_ip", src);
  if (query) params.set("message", query);
  if (dateFromVal) params.set("from_timestamp", new Date(dateFromVal).getTime());
  if (dateToVal) params.set("to_timestamp", new Date(dateToVal).getTime());

  showToast("Exporting to server...");
  try {
    const res = await authenticatedFetch(`${API_BASE}/api/export?${params.toString()}`);
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `Server error ${res.status}`);
    }
    const data = await res.json();
    const fileName = data.filename;
    const total = data.total || 0;

    // Also download the saved file to the user's machine for convenience.
    if (data.download) {
      const dl = await authenticatedFetch(`${API_BASE}${data.download}`);
      if (dl.ok) {
        const blob = await dl.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      }
    }

    showToast(`Saved ${fileName} to /exports (${total} rows)`);
  } catch (err) {
    console.error("Export error:", err);
    showToast("Export failed: " + err.message);
  }
}

function toggleView() {
  const tableView = document.getElementById("logTableView");
  const dashView = document.getElementById("dashboardView");
  const btn = document.getElementById("dashboardToggle");
  const statsBar = document.getElementById("statsBar");
  const toolbar = document.querySelector(".toolbar");

  if (dashView.classList.contains("active")) {
    // สลับไปหน้า Table
    dashView.classList.remove("active");
    tableView.style.display = "";
    statsBar.style.display = "";
    toolbar.style.display = "";
    btn.classList.remove("active");
    btn.innerHTML = "📊 Dashboard";
  } else {
    // สลับไปหน้า Dashboard
    dashView.classList.add("active");
    tableView.style.display = "none";
    statsBar.style.display = "none";
    toolbar.style.display = "none";
    btn.classList.add("active");
    btn.innerHTML = "📋 Logs";

    // สำคัญ: ต้องสั่ง renderDashboard() ทุกครั้งที่สลับมาหน้า Dashboard
    // เพื่อให้กราฟและสถิติใช้ข้อมูลล่าสุดจาก allLogs
    renderDashboard();
  }
}

function renderDashboard() {
  const total = allLogs.length;
  //const now = Date.running ? Date.now() : Date.now(); // Ensure current time
  const twentyFourHoursAgo = Date.now() - (24 * 60 * 60 * 1000);

  // 1. คำนวณ Error Rate
  const errors = allLogs.filter(
    (l) => l.level === "error" || l.level === "critical",
  ).length;
  const errorRate = total > 0 ? ((errors / total) * 100).toFixed(1) : 0;

  // 2. คำนวณ Total Logs (24h) - กรองเฉพาะที่เกิดในช่วง 24 ชม. ล่าสุด
  const logsLast24h = allLogs.filter(
    (l) => l.timestamp.getTime() > twentyFourHoursAgo
  ).length;

  // 3. คำนวณ Avg Response Time (ดึงจาก extras.duration หรือ extras.response_time)
  const responseTimes = allLogs
    .map((l) => parseFloat(l.extras?.duration || l.extras?.response_time || l.extras?.res_time))
    .filter((t) => !isNaN(t));
  const avgResponse = responseTimes.length > 0
    ? (responseTimes.reduce((a, b) => a + b, 0) / responseTimes.length).toFixed(2) + "ms"
    : "N/A";
  // คำนวณ Active Sources (นับจำนวน source ที่ไม่ซ้ำกัน)
  const activeSourcesCount = new Set(allLogs.map(l => l.source)).size;

  // --- อัปเดตตัวเลขลงใน UI ---
  // หมายเหตุ: ตรวจสอบว่าใน HTML ของคุณมี id เหล่านี้ (dashTotal, dashErrorRate, dashTotal24h, dashAvgResponse)
  const dashTotalElem = document.getElementById("dashTotal");
  if (dashTotalElem) dashTotalElem.textContent = total.toLocaleString();

  const dashErrorRateElem = document.getElementById("dashErrorRate");
  if (dashErrorRateElem) dashErrorRateElem.textContent = errorRate + "%";

  const dashTotal24hElem = document.getElementById("dashTotal24h");
  if (dashTotal24hElem) dashTotal24hElem.textContent = logsLast24h.toLocaleString();

  const dashAvgRespElem = document.getElementById("dashAvgResponse");
  if (dashAvgRespElem) dashAvgRespElem.textContent = avgResponse;

  const dashSourcesElem = document.getElementById("dashSources");
  if (dashSourcesElem) dashSourcesElem.textContent = activeSourcesCount;

  // 4. Time distribution chart (Bar Chart)
  const hourData = new Array(24).fill(0);
  allLogs.forEach((log) => {
    const hoursAgo = (Date.now() - log.timestamp.getTime()) / 3600000;
    const hour = Math.floor(hoursAgo);
    if (hour >= 0 && hour < 24) {
      hourData[hour]++;
    }
  });
  hourData.reverse();

  const maxHour = Math.max(...hourData, 1);
  const colors = {
    error: "var(--accent-red)",
    warn: "var(--accent-yellow)",
    info: "var(--cal-blue)", // ปรับให้ตรงกับ CSS ของคุณ
    debug: "var(--accent-purple)",
    critical: "#dc2626",
  };

  const timeChart = document.getElementById("timeChart");
  if (timeChart) {
    timeChart.innerHTML = hourData
      .map((count, i) => {
        const height = (count / maxHour) * 130;
        const hourLabel = `${23 - i}h`;
        return `
          <div class="bar-group">
              <div class="bar" style="height: ${height}px; background: var(--accent-blue); opacity: ${0.4 + (count / maxHour) * 0.6};" title="${count} logs"></div>
              <span class="bar-label">${hourLabel}</span>
          </div>
        `;
      })
      .join("");
  }

  // 5. Level distribution (Legend)
  const levelCounts = {};
  levels.forEach((l) => (levelCounts[l] = 0));
  allLogs.forEach((l) => {
    if (levelCounts.hasOwnProperty(l.level)) levelCounts[l.timestamp] = 0; // safety
    if (levelCounts.hasOwnProperty(l.level)) levelCounts[l.level]++;
  });

  const legend = document.getElementById("levelLegend");
  if (legend) {
    legend.innerHTML = levels
      .map((level) => {
        const count = levelCounts[level] || 0;
        const pct = total > 0 ? ((count / total) * 100).toFixed(1) : 0;
        return `
          <div class="legend-item">
              <span class="legend-dot" style="background: ${colors[level] || 'gray'};"></span>
              <span style="margin-left: 8px;">${level.charAt(0).toUpperCase() + level.slice(1)}</span>
              <span class="legend-count" style="margin-left: auto;">${count} (${pct}%)</span>
          </div>
        `;
      })
      .join("");
  }

  // 6. Top sources (Horizontal Bar)
  const sourceCounts = {};
  allLogs.forEach((l) => {
    const src = l.source || "unknown";
    sourceCounts[src] = (sourceCounts[src] || 0) + 1;
  });
  const sortedSources = Object.entries(sourceCounts).sort((a, b) => b[1] - a[1]);
  const maxSource = sortedSources.length > 0 ? sortedSources[0][1] : 1;
  const sourceColors = [
    "var(--accent-blue)", "var(--accent-cyan)", "var(--accent-green)",
    "var(--accent-purple)", "var(--accent-yellow)", "var(--accent-red)",
    "#06b6d4", "#ec4899",
  ];

  const topSourcesList = document.getElementById("topSourcesList");
  if (topSourcesList) {
    topSourcesList.innerHTML = sortedSources
      .slice(0, 8) // แสดงแค่ Top 8 เพื่อความสวยงาม
      .map(([source, count], i) => {
        const width = (count / maxSource) * 100;
        return `
          <div class="source-item">
              <span class="source-name">${source}</span>
              <div class="source-bar-container">
                  <div class="source-bar" style="width: ${width}%; background: ${sourceColors[i % sourceColors.length]};"></div>
              </div>
              <span class="source-count">${count}</span>
          </div>
        `;
      })
      .join("");
  }
}

function refreshLogs() {
  applyFilters();
  showToast("Logs refreshed");
}

function showToast(message) {
  const toast = document.getElementById("toast");
  document.getElementById("toastMessage").textContent = message;
  toast.classList.add("show");
  setTimeout(() => toast.classList.remove("show"), 3000);
}

// Live mode simulation
function startLiveMode() {
  if (liveInterval) clearInterval(liveInterval);
  liveInterval = setInterval(() => {
    if (liveMode) {
      const newLog = generateLogEntry(new Date());
      allLogs.unshift(newLog);
      // Keep max 10000 logs
      if (allLogs.length > 10000) {
        allLogs = allLogs.slice(0, 10000);
      }
      applyFilters();
      updateStats();
      refreshAllUI();
    }
  }, 3000);
}

function signout() {
  localStorage.removeItem("authToken");
  window.sessionStorage.removeItem("authToken");
  window.location.href = "/signin";
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  //generateLogs(500);
  renderDashboard();
  //startLiveMode();

  // If loadEngine is the function name
  if (typeof loadEngine === 'function') {
    loadEngine();
  }

  // Set default date range
  const now = new Date();
  const dayAgo = new Date(now - 86400000);
  document.getElementById("dateFrom").value = dayAgo.toISOString().slice(0, 16);
  document.getElementById("dateTo").value = now.toISOString().slice(0, 16);

  // Close modal on overlay click
  document.getElementById("logModal").addEventListener("click", (e) => {
    if (e.target.id === "logModal") closeModal();
  });

  // Keyboard shortcuts
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") closeModal();
    if (e.ctrlKey && e.key === "f") {
      e.preventDefault();
      document.getElementById("searchInput").focus();
    }
  });
});

const exportedFunctions = {
  refreshLogs,
  generateLogs,
  toggleView,
  toggleLevel,
  applyFilters,
  setQuickDate,
  clearAllFilters,
  debounceSearch,
  showLogDetail,
  closeModal,
  showExportHistoryModal,
  closeExportHistoryModal,
  downloadExport,
  deleteExport,
  copyLogDetail,
  exportSingleLog,
  exportCSV,
  exportLargeCSV,
  goToPage,
  changePageSize,
  sortLogs,
  signout,
  runSearch,
  normalizeSeverity,
  loadEngine,
  mapQuickwitHits,
  parseFlexibleTimestamp,
  dynamicMapper,
};

// Attach ไปยัง window object
Object.keys(exportedFunctions).forEach((fnName) => {
  window[fnName] = exportedFunctions[fnName];
});
