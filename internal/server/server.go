package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"distributed-URL-health-checker/internal/checker"
)

type Server struct {
	checker *checker.HealthChecker
	server  *http.Server
}

func New(hc *checker.HealthChecker, port string) *Server {
	s := &Server{
		checker: hc,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/health", s.handleHealth)

	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	statuses := s.checker.GetAllStatuses()

	html := `<!DOCTYPE html>
	<html>
	<head>
		<title>Health Checker Dashboard</title>
    <meta http-equiv="refresh" content="5">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        h1 { color: #333; }
        .status-grid { display: grid; gap: 20px; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); }
        .status-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .status-up { border-left: 4px solid #4CAF50; }
        .status-down { border-left: 4px solid #f44336; }
        .status-degraded { border-left: 4px solid #ff9800; }
        .url { font-weight: bold; margin-bottom: 10px; word-break: break-all; }
        .metric { margin: 5px 0; color: #666; }
        .error { color: #f44336; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>🔍 URL Health Checker Dashboard</h1>
    <p>Auto-refreshes every 5 seconds</p>
    <div class="status-grid">`

	for _, status := range statuses {
		html += fmt.Sprintf(`
        <div class="status-card status-%s">
            <div class="url">%s</div>
            <div class="metric">Status: <strong>%s</strong></div>
            <div class="metric">Response Time: %v</div>
            <div class="metric">Last Check: %s</div>
            %s
        </div>`,
			status.Status,
			status.URL,
			status.Status,
			status.ResponseTime.Round(time.Millisecond),
			status.LastCheck.Format("15:04:05"),
			func() string {
				if status.Error != "" {
					return fmt.Sprintf(`<div class="error">Error: %s</div>`, status.Error)
				}
				return ""
			}(),
		)
	}

	html += `
    </div>
</body>
</html>
	`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	statuses := s.checker.GetAllStatuses()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"statuses": statuses,
		"count":    len(statuses),
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}
