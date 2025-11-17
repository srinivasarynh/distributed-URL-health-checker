# Distributed URL Health Checker

A concurrent web service that monitors website availability in real-time.

## Features

- ✅ Monitors multiple URLs concurrently
- 📊 Live dashboard with auto-refresh
- 🚀 RESTful API
- 🔒 Thread-safe caching
- 🛡️ Graceful shutdown
- 🏎️ Production-grade patterns

## Quick Start

### 1. Setup

```bash
# Create project structure
mkdir -p healthchecker/cmd/healthchecker
mkdir -p healthchecker/internal/checker
mkdir -p healthchecker/internal/server
mkdir -p healthchecker/pkg/config
cd healthchecker

# Initialize module
go mod init healthchecker
```

### 2. Copy the code files to their locations:
```
healthchecker/
├── cmd/healthchecker/main.go
├── internal/
│   ├── checker/checker.go
│   ├── checker/cache.go
│   └── server/server.go
└── pkg/config/config.go
```

### 3. Run

```bash
go run cmd/healthchecker/main.go
```

Open browser: `http://localhost:8080`

## Configuration

Set environment variables:

```bash
# Monitor custom URLs
URLS="https://google.com,https://github.com,https://example.com" \
PORT=8080 \
CHECK_INTERVAL=10s \
TIMEOUT=5s \
go run cmd/healthchecker/main.go
```

### Options

| Variable | Default | Description |
|----------|---------|-------------|
| `URLS` | google.com, github.com, golang.org | Comma-separated URLs to monitor |
| `PORT` | 8080 | Server port |
| `CHECK_INTERVAL` | 10s | How often to check URLs |
| `TIMEOUT` | 5s | HTTP request timeout |

## API Endpoints

- `GET /` - Dashboard (HTML)
- `GET /api/status` - All URL statuses (JSON)
- `GET /api/health` - Server health (JSON)

### Example API Call

```bash
curl http://localhost:8080/api/status | jq
```

## Test with Race Detector

```bash
go run -race cmd/healthchecker/main.go
```

## Build for Production

```bash
# Build binary
go build -o healthchecker cmd/healthchecker/main.go

# Run
./healthchecker
```

## What You'll Learn

This project demonstrates:
- **Goroutines** - Concurrent execution
- **Channels** - Communication between goroutines
- **Worker Pools** - Controlled concurrency
- **Context** - Cancellation and timeouts
- **sync.Mutex** - Thread-safe shared data
- **sync.RWMutex** - Optimized read-heavy locking
- **sync.Once** - Lazy initialization
- **Package Structure** - Production organization

## Example Output

```
2024/11/17 10:30:00 Server starting on port 8080
```

Dashboard shows:
```
🔍 URL Health Checker Dashboard

┌─────────────────────────────┐
│ https://google.com          │
│ Status: up                  │
│ Response Time: 123ms        │
│ Last Check: 10:30:05        │
└─────────────────────────────┘
```

## Graceful Shutdown

Press `Ctrl+C`:
```
^C
2024/11/17 10:35:00 Shutting down gracefully...
2024/11/17 10:35:00 Server stopped
```

## License

MIT
