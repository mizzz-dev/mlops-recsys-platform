resource "google_cloud_run_v2_service" "inference_api" {
  name     = var.service_name
  location = var.region

  template {
    containers {
      image = var.image

      env {
        name  = "MODEL_PATH"
        value = var.model_path
      }

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}

resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  count    = var.allow_unauthenticated ? 1 : 0
  project  = var.project_id
  location = google_cloud_run_v2_service.inference_api.location
  name     = google_cloud_run_v2_service.inference_api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
