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

resource "clearblade_iot_device" "advanced-device" {
  registry = "registry-id"
  project  = "api-project-id"
  region   = "cloud-region"
  device = {
    id = "iot-gateway-101"
    gateway_config = {
      gateway_type        = "GATEWAY"
      gateway_auth_method = "ASSOCIATION_ONLY"
    }
  }
}
