
# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "ClearBlade/clearblade"
      version = "0.0.0-beta.6"
    }
  }
}

provider "clearblade" {
  
}

