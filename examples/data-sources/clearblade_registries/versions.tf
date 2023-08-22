terraform {
  required_providers {
    clearblade = {
      source = "ClearBlade/clearblade"
      version = "x.y.z" # check out the latest version in the release section
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = local.auth_credentials
  project     = local.gcp_project
  region      = local.gcp_region
}




