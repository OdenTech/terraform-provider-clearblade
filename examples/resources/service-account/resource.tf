# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account

resource "google_service_account" "iot_sa" {
  account_id   = "iot-sa"
  display_name = "IoT Service Account"
}

# note this requires the terraform to be run regularly
resource "time_rotating" "iot_sa_key_rotation" {
  rotation_days = 30
}

resource "google_service_account_key" "iot_sa_key" {
  service_account_id = google_service_account.iot_sa.name

  keepers = {
    rotation_time = time_rotating.iot_sa_key_rotation.rotation_rfc3339
  }
}

resource "google_project_iam_member" "iot_editor" {
  project = var.project_id
  role    = "roles/cloudiot.editor"
  member  = "serviceAccount:${google_service_account.iot_sa.email}"

  condition {
    title       = "expires_after_2022_07_31"
    description = "Expiring at midnight of 2022-07-31"
    expression  = "request.time < timestamp(\"2022-08-01T00:00:00Z\")"
  }
}

resource "google_project_iam_member" "pub_sub_editor" {
  project = var.project_id
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${google_service_account.iot_sa.email}"

  condition {
    title       = "expires_after_2022_07_31"
    description = "Expiring at midnight of 2022-07-31"
    expression  = "request.time < timestamp(\"2022-08-01T00:00:00Z\")"
  }
}

variable "project_id" {
  type = string
  default = "<project-id>"
}

output "iot_sa_private_key" {
  description = "Private key of the IoT service account"
  value       = google_service_account_key.iot_sa_key.private_key
  sensitive = true
}