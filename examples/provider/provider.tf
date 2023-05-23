
# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "clearblade.com/dev/clearblade"
    }
  }
}

provider "clearblade" {
  credentials = "credentials.json"
}

data "clearblade_registries" "all" {

}
