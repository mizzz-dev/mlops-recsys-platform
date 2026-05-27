# Staging smoke test

## Purpose

This document describes how to verify the Cloud Run staging service after deployment.

The smoke test is manual by design. It should be run after the staging deploy workflow completes.

## Workflow

Use the GitHub Actions workflow:

```text
Smoke staging
```

## Inputs

| Name | Required | Default | Description |
|---|---:|---|---|
| `service_url` | Yes | - | Cloud Run staging service URL |
| `use_identity_token` | No | `true` | Use a gcloud identity token for private Cloud Run services |
| `run_k6` | No | `true` | Run k6 smoke load test after endpoint checks |

## Checks

The workflow verifies:

- `/healthz`
- `/readyz`
- `/v1/recommendations?user_id=user_001&limit=3`
- JSON response validity
- `strategy` field presence
- `recommendations` field presence
- k6 smoke test, when enabled

## Authenticated Cloud Run

The staging deploy workflow uses `--no-allow-unauthenticated`, so authenticated requests are expected by default.

When `use_identity_token` is true, the workflow obtains an identity token using:

```bash
gcloud auth print-identity-token
```

The token is sent as:

```text
Authorization: Bearer <token>
```

## Public test endpoint

For a temporary public test endpoint, set:

```text
use_identity_token=false
```

Do not use this for production.

## Local equivalent

```bash
curl -s "$SERVICE_URL/healthz"
curl -s "$SERVICE_URL/readyz"
curl -s "$SERVICE_URL/v1/recommendations?user_id=user_001&limit=3"
make loadtest BASE_URL="$SERVICE_URL"
```

For authenticated k6:

```bash
ID_TOKEN=$(gcloud auth print-identity-token)
BASE_URL="$SERVICE_URL" ID_TOKEN="$ID_TOKEN" k6 run loadtests/k6/recommendation_api.js
```

## Pass criteria

- health check succeeds
- readiness check succeeds
- recommendation API returns valid JSON
- `strategy` is present
- `recommendations` is present
- k6 thresholds pass when enabled

## Failure handling

If the smoke test fails:

1. Check the failing step.
2. Review Cloud Run logs.
3. Check the deployed image URI.
4. Check model artifact availability.
5. Follow `docs/rollback.md` if user impact is possible.
