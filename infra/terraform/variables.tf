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

variable "allow_unauthenticated" {
  description = "Whether to allow unauthenticated invocation. Keep false for production."
  type        = bool
  default     = false
}
