variable "project_id" {
  description = "GCP project ID."
  type        = string
}

variable "region" {
  description = "GCP region for Cloud Run."
  type        = string
  default     = "asia-northeast1"
}

variable "service_name" {
  description = "Cloud Run service name."
  type        = string
  default     = "mlops-recsys-inference-api"
}

variable "image" {
  description = "Container image URI in Artifact Registry or another registry."
  type        = string
}

variable "model_path" {
  description = "Path to the model artifact inside the container."
  type        = string
  default     = "/app/artifacts/model.json"
}

variable "model_uri" {
  description = "Optional model artifact URI for future startup download. Example: gs://bucket/models/model.json."
  type        = string
  default     = ""
}

variable "model_artifact_bucket_name" {
  description = "GCS bucket name for model artifacts. Leave empty to use the generated default name."
  type        = string
  default     = ""
}

variable "cloud_run_runtime_service_account" {
  description = "Cloud Run runtime service account email. If empty, Cloud Run uses the platform default."
  type        = string
  default     = ""
}

variable "allow_unauthenticated" {
  description = "Whether to allow unauthenticated invocation. Keep false for production."
  type        = bool
  default     = false
}
