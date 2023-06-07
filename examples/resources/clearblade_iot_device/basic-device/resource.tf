# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "clearblade/clearblade"
    }
  }
}

provider "clearblade" {

}

# https://registry.terraform.io/providers/hashicorp/clearblade/latest/docs/resources/iot_device
 resource "clearblade_iot_device" "basic-device" {
  registry = "registry-id"
  project = "project-id"
  region  = "cloud-region"
  device = {
    id = "basic-device-100"
  }
}