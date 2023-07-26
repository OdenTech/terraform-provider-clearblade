# Service credential based configuration for the Clearblade IoT Core provider
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

# List all registries
data "clearblade_registries" "iot" {
  project = var.gcp_project_id
  region = var.gcp_region
}

output "iot_registries" {
  value = data.clearblade_registries.iot
}