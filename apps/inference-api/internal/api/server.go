package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"mlops-recsys-platform/apps/inference-api/internal/recommender"
)

type Server struct {
	svc       *recommender.Service
	requests  atomic.Uint64
	errors    atomic.Uint64
	fallbacks atomic.Uint64
	logger    *slog.Logger
}

type RecommendationResponse struct {
	UserID          string                       `json:"user_id"`
	ModelVersion    *string                      `json:"model_version"`
	Strategy        string                       `json:"strategy"`
	Recommendations []recommender.Recommendation `json:"recommendations"`
	LatencyMS       int64                        `json:"latency_ms"`
}

func NewServer(svc *recommender.Service, logger *slog.Logger) *Server {
	return &Server{svc: svc, logger: logger}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("GET /metrics", s.metrics)
	mux.HandleFunc("GET /v1/models/current", s.currentModel)
	mux.HandleFunc("GET /v1/recommendations", s.recommendations)
	return mux
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"model_loaded":  s.svc.ModelLoaded(),
		"model_version": s.svc.ModelVersion(),
		"fallback":      true,
	})
}

func (s *Server) currentModel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"model_loaded":  s.svc.ModelLoaded(),
		"model_version": s.svc.ModelVersion(),
	})
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	modelLoaded := 0
	if s.svc.ModelLoaded() {
		modelLoaded = 1
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP recsys_requests_total Total recommendation API requests.\n")
	fmt.Fprintf(w, "# TYPE recsys_requests_total counter\n")
	fmt.Fprintf(w, "recsys_requests_total %d\n", s.requests.Load())
	fmt.Fprintf(w, "# HELP recsys_errors_total Total recommendation API errors.\n")
	fmt.Fprintf(w, "# TYPE recsys_errors_total counter\n")
	fmt.Fprintf(w, "recsys_errors_total %d\n", s.errors.Load())
	fmt.Fprintf(w, "# HELP recsys_fallback_total Total fallback recommendation responses.\n")
	fmt.Fprintf(w, "# TYPE recsys_fallback_total counter\n")
	fmt.Fprintf(w, "recsys_fallback_total %d\n", s.fallbacks.Load())
	fmt.Fprintf(w, "# HELP recsys_model_loaded Whether model artifact is loaded.\n")
	fmt.Fprintf(w, "# TYPE recsys_model_loaded gauge\n")
	fmt.Fprintf(w, "recsys_model_loaded %d\n", modelLoaded)
}

func (s *Server) recommendations(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	s.requests.Add(1)

	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		s.errors.Add(1)
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		s.errors.Add(1)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := s.svc.Recommend(userID, limit)
	if result.Strategy == recommender.StrategyFallback {
		s.fallbacks.Add(1)
	}

	writeJSON(w, http.StatusOK, RecommendationResponse{
		UserID:          userID,
		ModelVersion:    result.ModelVersion,
		Strategy:        result.Strategy,
		Recommendations: result.Recommendations,
		LatencyMS:       time.Since(started).Milliseconds(),
	})
}

func parseLimit(value string) (int, error) {
	if value == "" {
		return 5, nil
	}
	limit, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("limit must be an integer")
	}
	if limit < 1 || limit > 20 {
		return 0, fmt.Errorf("limit must be between 1 and 20")
	}
	return limit, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
