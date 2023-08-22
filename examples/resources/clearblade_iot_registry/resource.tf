terraform {
  required_providers {
    clearblade = {
      source  = "ClearBlade/clearblade"
      version = "0.2.4"
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = var.clearblade-creds
  project     = var.gcp_project_id
  region      = var.gcp_region
}

resource "clearblade_iot_registry" "example" {
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
    mqtt_enabled_state = "MQTT_ENABLED"
  }

  http_config = {
    http_enabled_state = "HTTP_ENABLED"
  }

  log_level = var.log_level
}

resource "clearblade_iot_registry" "example1" {
  id = var.registry_id_1

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
    mqtt_enabled_state = "MQTT_ENABLED"
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  log_level = var.log_level
}


