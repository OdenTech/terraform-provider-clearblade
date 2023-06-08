
# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "clearblade/clearblade"
      version = "0.0.0-beta.4"
    }
  }
}

provider "clearblade" {
  
}

data "clearblade_registries" "all" {

}
