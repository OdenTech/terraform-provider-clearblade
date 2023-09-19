# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source  = "ClearBlade/clearblade"
      version = "x.y.z" # check out the latest version in the release section
    }
  }
}
# List all registries
data "clearblade_registries" "example" {

}

output "all_registries" {
  value = data.clearblade_registries.example
}
