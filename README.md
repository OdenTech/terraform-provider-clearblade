# Clearblade Terraform Provider

## Getting started

Add this provider in your terraform configuration block:

```terraform
terraform {
  required_providers {
    clearblade = {
      source = "clearblade/clearblade"
      version = "0.0.0-beta.3"
    }
  }
}

# Configure the planetscale provider
provider "clearblade" {

}
```
