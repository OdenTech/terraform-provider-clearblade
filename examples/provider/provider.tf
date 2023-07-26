
# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source = "ClearBlade/clearblade"
      version = "0.0.0-beta.7"
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = var.clearblade-creds
}

