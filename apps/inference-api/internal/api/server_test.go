package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"mlops-recsys-platform/apps/inference-api/internal/recommender"
)

func TestRecommendationsReturnsModelResult(t *testing.T) {
	model := &recommender.Model{Version: "test-model", Recommendations: []recommender.Recommendation{{ContentID: "quest_101", Score: 1, Reason: "popular_content"}}}
	server := NewServer(recommender.NewService(model), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/v1/recommendations?user_id=user_001&limit=1", nil)
	res := httptest.NewRecorder()
	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	var body RecommendationResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Strategy != recommender.StrategyModel {
		t.Fatalf("expected model strategy, got %s", body.Strategy)
	}
	if body.ModelVersion == nil || *body.ModelVersion != "test-model" {
		t.Fatalf("unexpected model version: %#v", body.ModelVersion)
	}
}

func TestRecommendationsRequiresUserID(t *testing.T) {
	server := NewServer(recommender.NewService(nil), slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/v1/recommendations", nil)
	res := httptest.NewRecorder()
	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.Code)
	}
}

func TestRecommendationsValidatesLimit(t *testing.T) {
	server := NewServer(recommender.NewService(nil), slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/v1/recommendations?user_id=user_001&limit=100", nil)
	res := httptest.NewRecorder()
	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.Code)
	}
}

func TestRecommendationsFallsBackWhenModelMissing(t *testing.T) {
	server := NewServer(recommender.NewService(nil), slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/v1/recommendations?user_id=user_001&limit=2", nil)
	res := httptest.NewRecorder()
	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	var body RecommendationResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Strategy != recommender.StrategyFallback {
		t.Fatalf("expected fallback strategy, got %s", body.Strategy)
	}
	if len(body.Recommendations) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(body.Recommendations))
	}
}
