locals {
  model_artifact_bucket_name = var.model_artifact_bucket_name != "" ? var.model_artifact_bucket_name : "${var.project_id}-mlops-recsys-model-artifacts"
}

resource "google_storage_bucket" "model_artifacts" {
  name                        = local.model_artifact_bucket_name
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      num_newer_versions = 10
    }
    action {
      type = "Delete"
    }
  }
}

resource "google_storage_bucket_iam_member" "model_artifact_reader" {
  count  = var.cloud_run_runtime_service_account == "" ? 0 : 1
  bucket = google_storage_bucket.model_artifacts.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${var.cloud_run_runtime_service_account}"
}
