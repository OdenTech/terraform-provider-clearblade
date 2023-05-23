terraform {
  required_providers {
    clearblade = {
      version = "1.0.0"
      source  = "clearblade/clearblade"
    }
  }
}

provider "clearblade" {}

data "clearblade_registry" "example" {}
