# Service credential based configuration for the Clearblade IoT Core provider
terraform {
  required_providers {
    clearblade = {
      source  = "ClearBlade/clearblade"
      version = "0.2.3"
    }
  }
}
# List all registries
data "clearblade_registries" "example" {

}

output "all_registries" {
  value = data.clearblade_registries.example
}
