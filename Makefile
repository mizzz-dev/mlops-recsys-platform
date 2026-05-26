SHELL := /bin/bash
API_DIR := apps/inference-api
PYTHON := python3

.PHONY: help test go-test py-test train run-api request-sample compile-pipeline docker-build clean

help:
	@echo "Available commands:"
	@echo "  make test              Run Go and Python tests"
	@echo "  make train             Generate synthetic data and train local model"
	@echo "  make run-api           Run inference API locally"
	@echo "  make request-sample    Call recommendation endpoint"
	@echo "  make compile-pipeline  Compile local pipeline spec"
	@echo "  make docker-build      Build inference API Docker image"
	@echo "  make clean             Remove generated local artifacts"

test: go-test py-test compile-pipeline

go-test:
	cd $(API_DIR) && go test ./...

py-test:
	PYTHONPATH=ml/src $(PYTHON) -m pytest ml/tests

train:
	PYTHONPATH=ml/src $(PYTHON) -m mlops_recsys.train --events data/samples/events.csv --model artifacts/model.json

run-api:
	cd $(API_DIR) && MODEL_PATH=../../artifacts/model.json go run ./cmd/server

request-sample:
	curl -s 'http://localhost:8080/v1/recommendations?user_id=user_001&limit=3' | $(PYTHON) -m json.tool

compile-pipeline:
	$(PYTHON) pipelines/training/pipeline.py --compile --output artifacts/pipeline.json

docker-build:
	docker build -t mlops-recsys-inference-api:local $(API_DIR)

clean:
	rm -f artifacts/model.json artifacts/pipeline.json data/samples/events.csv data/samples/features.json
