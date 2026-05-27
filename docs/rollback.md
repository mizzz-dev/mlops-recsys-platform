# Rollback guide

## Purpose

This document describes how to roll back a staging Cloud Run deployment when the recommendation API has degraded.

## When to roll back

Consider rollback when one or more of the following happens after deployment:

- `/healthz` does not return 200
- `/readyz` does not return a valid response
- p95 latency exceeds the agreed threshold
- error rate exceeds 1%
- fallback rate unexpectedly increases
- model version is not the expected version

## Pre-check

```bash
gcloud run revisions list \
  --service mlops-recsys-inference-api-staging \
  --region asia-northeast1
```

Identify the last known good revision.

## Roll back all traffic to a previous revision

```bash
gcloud run services update-traffic mlops-recsys-inference-api-staging \
  --region asia-northeast1 \
  --to-revisions <REVISION_NAME>=100
```

## Verify after rollback

```bash
SERVICE_URL=$(gcloud run services describe mlops-recsys-inference-api-staging \
  --region asia-northeast1 \
  --format 'value(status.url)')

curl -s "$SERVICE_URL/healthz"
curl -s "$SERVICE_URL/readyz"
curl -s "$SERVICE_URL/v1/recommendations?user_id=user_001&limit=3"
make loadtest BASE_URL="$SERVICE_URL"
```

## Incident note template

Record the following after rollback:

- Date and time
- Deployed revision
- Rolled back revision
- Trigger metric or symptom
- User impact
- Root cause hypothesis
- Follow-up issue

## Notes

This project keeps production deployment out of scope for now. Production rollback should require stricter approval, audit trail, and communication steps.
