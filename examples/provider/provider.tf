
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

data "clearblade_registries" "all" {

}
