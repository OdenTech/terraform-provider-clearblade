
# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source  = "ClearBlade/clearblade"
      version = "0.1.1"
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = var.clearblade-creds
  project     = var.gcp_project_id
  region      = var.gcp_region
}

