# Model artifact storage design

## Purpose

This document describes the target design for storing recommendation model artifacts outside the inference API container.

The current MVP reads a local file from `MODEL_PATH`. The next production-oriented step is to store model artifacts in GCS and let the serving layer resolve the artifact through `MODEL_URI`.

## Current MVP behavior

```text
make train
  -> artifacts/model.json

inference-api
  -> reads MODEL_PATH
  -> serves recommendations
  -> falls back when the model file is missing
```

## Target behavior

```text
training pipeline
  -> writes model artifact to GCS
  -> records model URI and evaluation metadata

Cloud Run inference API
  -> receives MODEL_URI=gs://<bucket>/models/<version>/model.json
  -> downloads artifact at startup or during a controlled reload
  -> validates model schema
  -> serves recommendations
  -> falls back when the model cannot be loaded
```

## Environment variables

| Name | Purpose | Status |
|---|---|---|
| `MODEL_PATH` | Local path used by the current API implementation | Implemented |
| `MODEL_URI` | Remote model artifact URI such as `gs://bucket/models/model.json` | Terraform only |

## GCS layout

Recommended object layout:

```text
gs://<bucket>/models/<model-version>/model.json
gs://<bucket>/models/<model-version>/metrics.json
gs://<bucket>/models/latest.json
```

Example:

```text
gs://my-project-mlops-recsys-model-artifacts/models/popular-baseline-local/model.json
```

## IAM

Cloud Run runtime service account needs read-only access to the model artifact bucket.

Required role:

```text
roles/storage.objectViewer
```

Terraform supports this through:

```text
cloud_run_runtime_service_account
```

When the variable is set, Terraform grants `roles/storage.objectViewer` on the model artifact bucket.

## Model versioning

Model artifacts should be immutable once published.

Recommended policy:

- write artifacts under a versioned prefix
- keep GCS bucket versioning enabled
- use `latest.json` only as a pointer
- do not overwrite historical model versions
- retain enough historical versions for rollback

## Failure behavior

The API must not fail hard when the remote artifact cannot be loaded.

Expected behavior:

- log model download failure
- keep serving fallback recommendations
- expose model state through `/readyz`
- expose fallback count through `/metrics`

## Non-goals in this PR

- API startup download from GCS
- model registry implementation
- signed URL support
- automatic model reload
- training pipeline upload implementation

## Follow-up implementation tasks

1. Add `MODEL_URI` parsing in the Go inference API.
2. Add GCS downloader behind a small interface.
3. Download model artifact to a temporary local path during startup.
4. Validate model JSON before replacing the active model.
5. Add tests for invalid URI, download failure, and fallback behavior.
6. Add training pipeline upload step.
7. Add model metadata file such as `metrics.json` or `latest.json`.
