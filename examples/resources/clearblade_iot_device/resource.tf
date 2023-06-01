# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "clearblade.com/dev/clearblade"
    }
  }
}

provider "clearblade" {
  
}

resource "clearblade_iot_device" "basic-device" {
  name = "basic-device-111"
  registry = "terraform-testregistry1"
  project  = "api-project-320446546234"
  region   = "us-central1"
  //registry = google_cloudiot_registry.iot-registry.id
}

//data "clearblade_registries" "edu" {}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_device

/* resource "google_cloudiot_device" "basic-device" {
  name     = "basic-device"
  registry = google_cloudiot_registry.iot-registry.id
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_device

resource "google_cloudiot_device" "advanced-device" {
  name     = "advanced-device"
  registry = google_cloudiot_registry.iot-registry.id

  credentials {
    public_key {
      format = "RSA_X509_PEM"
      key    = file("~/.auth/rsa_cert.pem")
    }
  }
} */
