# Artifact Registry image publishing

## Purpose

This document describes how to publish the inference API container image to Artifact Registry for staging deployment.

The workflow is intentionally manual and safe by default. It builds and pushes the image, but it does not deploy the service.

## Workflow

Use the GitHub Actions workflow:

```text
Publish inference API image
```

The workflow is triggered by `workflow_dispatch`.

## Required GitHub environment

Create or reuse the `staging` GitHub environment.

Recommended protection rules:

- Require manual approval
- Limit who can approve deployments
- Keep production separated from staging

## Required GitHub secrets

| Name | Purpose |
|---|---|
| `GCP_PROJECT_ID` | Target GCP project ID |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Workload Identity Provider resource name |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | Service account email used by GitHub Actions |

## Required GitHub variables

| Name | Default | Purpose |
|---|---|---|
| `GCP_REGION` | `asia-northeast1` | Artifact Registry region |
| `ARTIFACT_REGISTRY_REPOSITORY` | `mlops-recsys` | Artifact Registry repository name |

## Required GCP setup

Enable the API:

```bash
gcloud services enable artifactregistry.googleapis.com
```

Create a Docker repository:

```bash
gcloud artifacts repositories create mlops-recsys \
  --repository-format=docker \
  --location=asia-northeast1 \
  --description="MLOps recommendation platform images"
```

## Required IAM roles

Grant the GitHub deploy service account permission to push images.

Recommended role:

- `roles/artifactregistry.writer`

Avoid broad owner/editor roles.

## Image naming policy

The workflow publishes the image as:

```text
<region>-docker.pkg.dev/<project-id>/<repository>/mlops-recsys-inference-api:<tag>
```

Default tag:

```text
${GITHUB_SHA}
```

You can provide an explicit `image_tag` when running the workflow. Prefer immutable tags such as commit SHA or release candidate tags.

## How to deploy after publishing

Copy the printed image URI and pass it to the staging deploy workflow:

```text
Deploy staging
```

The deploy workflow accepts the image URI as the `image` input.

## Notes

- This workflow does not deploy automatically.
- This workflow does not publish production images.
- Do not store service account keys in GitHub.
- Use Workload Identity Federation.
