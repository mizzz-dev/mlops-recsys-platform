# Staging release guide

## Purpose

This guide connects the image publish workflow, staging deploy workflow, and staging smoke test workflow.

The release is intentionally manual. This keeps the portfolio safe by default and makes each operational step reviewable.

## Release flow

```text
1. Run CI on pull request
2. Merge to main
3. Run Publish inference API image
4. Copy the published image URI
5. Run Deploy staging with that image URI
6. Run Smoke staging with the deployed service URL
7. Roll back if the verification fails
```

## Step 1: Publish image

Open GitHub Actions and run:

```text
Publish inference API image
```

Inputs:

| Name | Required | Description |
|---|---:|---|
| `image_tag` | No | Optional tag. Defaults to the workflow commit SHA. |

After the workflow completes, copy the `image_uri` from the job summary.

Example format:

```text
asia-northeast1-docker.pkg.dev/<project-id>/<repository>/mlops-recsys-inference-api:<tag>
```

## Step 2: Deploy staging

Open GitHub Actions and run:

```text
Deploy staging
```

Inputs:

| Name | Required | Description |
|---|---:|---|
| `image` | Yes | Image URI copied from the publish workflow summary. |
| `service_name` | No | Cloud Run staging service name. |

## Step 3: Smoke test staging

Open GitHub Actions and run:

```text
Smoke staging
```

Inputs:

| Name | Required | Description |
|---|---:|---|
| `service_url` | Yes | Cloud Run staging service URL. |
| `use_identity_token` | No | Use identity token for private Cloud Run services. |
| `run_k6` | No | Run k6 smoke load test. |

For details, see:

```text
docs/staging-smoke-test.md
```

## Acceptance criteria

- `/healthz` returns 200
- `/readyz` returns 200 and model state is visible
- recommendation API returns `strategy` and `recommendations`
- smoke load test passes when enabled
- unexpected fallback increase is not observed

## Rollback

If verification fails, follow:

```text
docs/rollback.md
```

Use Cloud Run traffic rollback to restore the last known good revision.

## Notes

- Production deployment is out of scope.
- Service account keys must not be used.
- Use Workload Identity Federation.
- Keep image tags immutable, preferably commit SHA based.
