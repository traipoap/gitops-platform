exporter/
├── api/                   # API routes (ถ้ามี)
├── controllers/           # Logic สำหรับจัดการ request และ response
│   ├── export_controller.go # Controller สำหรับ Export Data
│   ├── search_controller.go # Controller สำหรับ Search Functionality
│   └── quickwit_controller.go # Controller สำหรับ QuickWIt Functionality
├── models/                # โดมินิเซน (Data Models)
│   ├── search.go          # Definition ของโครงสร้างข้อมูลสำหรับ search
│   └── export.go          # Definition ของโครงสร้างข้อมูลสำหรับ export
├── services/              # Business logic สำหรับการประมวลผลข้อมูล
│   ├── export_service.go  # Service สำหรับการ export data
│   └── quickwit_client.go # Service สำหรับการเชื่อมต่อกับ Quickwit
├── routers/               # กำหนดเส้นทาง (Routes)
│   └── routers.go         # จัดการ Routing logic
├── pages/                 # ไฟล์ UI (HTML, CSS, JS)
│   ├── ui1.html           # HTML สำหรับหน้า UI ที่ 1
│   └── ui2.html           # HTML สำหรับหน้า UI ที่ 2
├── static/                # Assets Static (CSS, JS, Images)
│   ├── css/              # CSS files
│   │   └── style.css     # Main CSS file
│   ├── js/               # JavaScript files
│   │   └── script.js     # Main JavaScript file
│   └── img/             # Image files
├── .gitignore           # ไฟล์ที่ต้อง ignore ใน Git
├── go.mod                # Module definition (Go Modules)
├── go.sum                # checksum สำหรับ Go modules
└── main.go               # Entry point ของ application
