# GCP setup for staging deployment

## Purpose

This document describes the minimum GCP and GitHub configuration required before enabling the staging deployment workflow.

The workflow is intentionally manual and safe by default. It does not build or push images by itself. The image URI must be provided explicitly when running the workflow.

## Required GitHub environment

Create a GitHub environment named `staging`.

Recommended protection rules:

- Require manual approval before deployment
- Limit who can approve deployments
- Keep production separate from staging

## Required GitHub secrets

| Name | Purpose |
|---|---|
| `GCP_PROJECT_ID` | Target GCP project ID |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Workload Identity Provider resource name |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | Service account email used by GitHub Actions |

## Optional GitHub variables

| Name | Default | Purpose |
|---|---|---|
| `GCP_REGION` | `asia-northeast1` | Cloud Run deployment region |

## Required GCP services

Enable the following APIs:

```bash
gcloud services enable run.googleapis.com
gcloud services enable iamcredentials.googleapis.com
gcloud services enable artifactregistry.googleapis.com
```

## Required IAM roles

Grant the deployment service account the minimum roles required for staging.

Recommended starting point:

- `roles/run.admin`
- `roles/iam.serviceAccountUser` for the Cloud Run runtime service account
- Artifact Registry read permission if the image is private

Avoid broad owner/editor roles.

## Manual workflow input

The staging workflow requires an explicit container image URI.

Example:

```text
asia-northeast1-docker.pkg.dev/<project-id>/<repository>/mlops-recsys-inference-api:<tag>
```

## Security notes

- Do not enable unauthenticated access by default.
- Keep production deployment in a separate workflow and environment.
- Use short-lived credentials through Workload Identity Federation.
- Do not store service account keys in GitHub secrets.
