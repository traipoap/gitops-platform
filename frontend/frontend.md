# Project Context: APP - System Log Manager

## 📋 Project Overview
APP is a self-hosted system log management dashboard for viewing, filtering, and analyzing server logs. Built with a focus on performance, simplicity, and developer experience.

**Key Features:**
- Real-time log streaming with live mode
- Advanced filtering (level, source, host, time range, search)
- Interactive dashboard with charts and metrics
- Export logs to CSV/JSON
- Dark theme UI optimized for long monitoring sessions
- Authentication-ready login/signup pages

## 🛠️ Tech Stack

### Core
- **Framework**: Astro (SSG/SSR hybrid)
- **Language**: TypeScript (preferred) / Vanilla JavaScript
- **Styling**: Plain CSS with CSS variables (no preprocessor)
- **Build Tool**: Vite (via Astro)

### UI/UX
- **Fonts**: Inter (UI), JetBrains Mono (code/logs)
- **Icons**: Inline SVG (no icon library dependency)
- **Theme**: Dark mode only (CSS variables in `:root`)
- **Animations**: CSS keyframes + minimal JS

### Architecture Pattern
- **Island Architecture**: Astro components with selective hydration
- **Progressive Enhancement**: Vanilla JS for interactivity, no heavy framework required
- **Module Pattern**: Functions exposed to `window` for inline `onclick` handlers (legacy compatibility)

## 📁 Project Structure
.
├── assets
│   ├── astro.svg
│   ├── background.svg
│   └── logout.svg
├── components
│   └── Welcome.astro
├── layouts
│   └── Layout.astro
├── pages
│   ├── dashboard.astro
│   ├── export.astro         # Export page for log data
│   ├── signin.astro
│   └── signup.astro
├── scripts
│   └── main.js
└── styles
    ├── main.css
    ├── signin.css
    └── signup.css

7 directories, 13 files

### API Backend
- localhost:8080/api/auth/login
- localhost:8080/task/search
- localhost:8080/task/indies