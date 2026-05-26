package recommender

import "strings"

const (
	StrategyModel    = "model"
	StrategyFallback = "fallback_popular"
)

type Result struct {
	ModelVersion    *string          `json:"model_version"`
	Strategy        string           `json:"strategy"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Service struct {
	model    *Model
	fallback []Recommendation
}

func NewService(model *Model) *Service {
	return &Service{
		model: model,
		fallback: []Recommendation{
			{ContentID: "quest_001", Score: 1.0, Reason: "fallback_popular"},
			{ContentID: "quest_002", Score: 0.9, Reason: "fallback_popular"},
			{ContentID: "quest_003", Score: 0.8, Reason: "fallback_popular"},
			{ContentID: "quest_004", Score: 0.7, Reason: "fallback_popular"},
			{ContentID: "quest_005", Score: 0.6, Reason: "fallback_popular"},
		},
	}
}

func (s *Service) ModelLoaded() bool {
	return s.model != nil
}

func (s *Service) ModelVersion() *string {
	if s.model == nil {
		return nil
	}
	version := s.model.Version
	return &version
}

func (s *Service) Recommend(userID string, limit int) Result {
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	if strings.TrimSpace(userID) == "" || s.model == nil {
		return Result{ModelVersion: nil, Strategy: StrategyFallback, Recommendations: top(s.fallback, limit)}
	}

	version := s.model.Version
	return Result{ModelVersion: &version, Strategy: StrategyModel, Recommendations: s.model.Top(limit)}
}

func top(items []Recommendation, limit int) []Recommendation {
	if limit > len(items) {
		limit = len(items)
	}
	out := make([]Recommendation, limit)
	copy(out, items[:limit])
	return out
}
