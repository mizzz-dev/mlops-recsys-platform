output "service_name" {
  description = "Cloud Run service name."
  value       = google_cloud_run_v2_service.inference_api.name
}

output "service_uri" {
  description = "Cloud Run service URI."
  value       = google_cloud_run_v2_service.inference_api.uri
}

output "model_artifact_bucket_name" {
  description = "GCS bucket name for model artifacts."
  value       = google_storage_bucket.model_artifacts.name
}

output "model_artifact_bucket_uri" {
  description = "GCS URI for model artifact bucket."
  value       = "gs://${google_storage_bucket.model_artifacts.name}"
}
