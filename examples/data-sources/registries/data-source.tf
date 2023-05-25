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

# List all registries
data "clearblade_registries" "all" {
  project = "api-project-320446546234"
  region  = "us-central1"
}
