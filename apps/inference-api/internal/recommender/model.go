package recommender

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

type Recommendation struct {
	ContentID string  `json:"content_id"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason"`
}

type Model struct {
	Version         string           `json:"version"`
	GeneratedAt     string           `json:"generated_at"`
	Evaluation      map[string]any   `json:"evaluation"`
	Recommendations []Recommendation `json:"recommendations"`
}

var ErrModelNotFound = errors.New("model artifact not found")

func LoadModel(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrModelNotFound
		}
		return nil, fmt.Errorf("open model artifact: %w", err)
	}
	defer file.Close()

	var model Model
	if err := json.NewDecoder(file).Decode(&model); err != nil {
		return nil, fmt.Errorf("decode model artifact: %w", err)
	}
	if model.Version == "" {
		return nil, errors.New("model version is empty")
	}
	if len(model.Recommendations) == 0 {
		return nil, errors.New("model recommendations are empty")
	}

	sort.SliceStable(model.Recommendations, func(i, j int) bool {
		return model.Recommendations[i].Score > model.Recommendations[j].Score
	})
	return &model, nil
}

func (m *Model) Top(limit int) []Recommendation {
	if limit > len(m.Recommendations) {
		limit = len(m.Recommendations)
	}
	items := make([]Recommendation, limit)
	copy(items, m.Recommendations[:limit])
	return items
}
