# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      version = "0.0.0-beta.4"
      source = "clearblade.com/dev/clearblade"
    }
  }
}

provider "clearblade" {

}
# List all registries
data "clearblade_registries" "all" {
  project = "project-id"
  region  = "us-central1"
}
