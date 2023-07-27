terraform {
  required_providers {
    clearblade = {
      source  = "ClearBlade/clearblade"
      version = "0.0.0-beta.7"
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = var.clearblade-creds
}

resource "clearblade_iot_registry" "example-registry" {
  project = var.gcp_project_id
  region  = var.gcp_region
  registry = {
    id = var.registry_id

    event_notification_configs = [
      {
        pubsub_topic_name = var.event_topic_name
        subfolder_matches = var.event_subfolder_matches
      },

      {
        pubsub_topic_name = var.event_topic_name
        subfolder_matches = ""
      }
    ]

    state_notification_config = {
      pubsub_topic_name = var.state_topic_name
    }

    mqtt_config = {
      mqtt_config = "MQTT_ENABLED"
    }

    http_config = {
      http_config = "HTTP_DISABLED"
    }

    log_level = var.log_level
  }

}
