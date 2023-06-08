# Clearblade Terraform Provider

A modern (protocol 6) Terraform provider for the Clearblade IoT Core service (work-in-progress).

## `Development status`

This Terraform provider code for the Clearblade IoT Core service is currently in Preview.

## Getting started

Add this provider in your terraform configuration block:

```terraform
terraform {
  required_providers {
    clearblade = {
      source = "clearblade/clearblade"
      version = "0.0.0-beta.4"
    }
  }
}

# Configure the Clearblade provider
provider "clearblade" {

}
```


