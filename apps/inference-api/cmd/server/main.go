package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"mlops-recsys-platform/apps/inference-api/internal/api"
	"mlops-recsys-platform/apps/inference-api/internal/recommender"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	modelPath := getenv("MODEL_PATH", "../../artifacts/model.json")
	port := getenv("PORT", "8080")

	model, err := recommender.LoadModel(modelPath)
	if err != nil {
		if errors.Is(err, recommender.ErrModelNotFound) {
			logger.Warn("model artifact not found; fallback recommendations will be used", "path", modelPath)
		} else {
			logger.Warn("model artifact could not be loaded; fallback recommendations will be used", "path", modelPath, "error", err)
		}
	} else {
		logger.Info("model artifact loaded", "path", modelPath, "version", model.Version)
	}

	svc := recommender.NewService(model)
	server := api.NewServer(svc, logger)
	logger.Info("starting inference api", "port", port)
	if err := http.ListenAndServe(":"+port, server.Routes()); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
