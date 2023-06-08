# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      version = "0.0.0-beta.3"
      source  = "clearblade.com/clearblade"
    }
  }
}

provider "clearblade" {

}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_registry

resource "clearblade_iot_registry" "iot-registry" {
  name = "iot-registry"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.additional-telemetry.id
    subfolder_matches = "test/path"
  }

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.default-telemetry.id
    subfolder_matches = ""
  }

  state_notification_config = {
    pubsub_topic_name = google_pubsub_topic.default-devicestatus.id
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_ENABLED"
  }

  http_config = {
    http_enabled_state = "HTTP_ENABLED"
  }

  log_level = "INFO"
}

resource "google_pubsub_topic" "default-devicestatus" {
  name = "default-devicestatus"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "default-telemetry"
}

resource "google_pubsub_topic" "additional-telemetry" {
  name = "additional-telemetry"
}
