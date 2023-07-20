# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      version = "0.0.0-beta.5"
      source  = "clearblade/clearblade"
    }
  }
}

provider "clearblade" {

}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_registry

resource "clearblade_iot_registry" "iot-registry" {
  project = "project-id"
  region  = "us-central1"
  registry = {
    id = "tf-registry-101"

    event_notification_configs = [
      {
        pubsub_topic_name = "projects/api-project-320446546234/topics/rootevent"
        subfolder_matches = "test/path"
      },

      {
        pubsub_topic_name = "projects/api-project-320446546234/topics/rootevent"
        subfolder_matches = ""
      }
    ]

    state_notification_config = {
      pubsub_topic_name = "projects/api-project-320446546234/topics/rootevent"
    }

    mqtt_config = {
      mqtt_enabled_state = "MQTT_ENABLED"
    }

    http_config = {
      http_enabled_state = "HTTP_ENABLED"
    }

    log_level = "INFO"
  }

}