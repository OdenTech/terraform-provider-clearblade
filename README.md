<div align="center">
   <p>ClearBlade Terraform provider (Beta).</p>
   <a href="https://clearblade.atlassian.net/wiki/spaces/IC/overview"><img src="https://img.shields.io/static/v1?label=Docs&message=API Ref&color=000000&style=for-the-badge" /></a>
  <a href="https://opensource.org/licenses/MPL-2.0"><img src="https://img.shields.io/badge/License-MPL-blue.svg?style=for-the-badge" /></a>
</div>

## Authentication

Developers will need to create or download a ClearBlade service account credential within your [ClearBlade IoT Core Developer Portal](https://iot.clearblade.com/iot-core/) to make API requests. You can use your existing ClearBlade IoT Core account to log in to the Developer Portal. Once you are in the Developer Portal, [add service accounts to a project](https://clearblade.atlassian.net/wiki/spaces/IC/pages/2240675843/) and download the JSON file with your service account's credentials.

Use the provider block to configure the path to your service account's JSON file for authentication. Optionally, the following environment variable can be used to set your credentials in your terminal or IDE environment:

```
 export CLEARBLADE_CONFIGURATION=/path/to/file.json
```

<!-- Start SDK Installation -->

## SDK Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    clearblade = {
      source = "ClearBlade/clearblade"
      version = "0.1.0"
    }
  }
}

provider "clearblade" {
  # Configuration options
  credentials = var.clearblade-creds
}
```

<!-- End SDK Installation -->

## Testing the provider locally

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/ClearBlade/clearblade" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Your `<PATH>` may vary depending on how your Go environment variables are configured. Execute `go env GOBIN` to set it, then set the `<PATH>` to the value returned. If nothing is returned, set it to the default location, `$HOME/go/bin`.

## Creating a new release

To create a new release execute the following commands:

```shell
# Use sem var for tags, i.e. v0.0.1
git tag [tag]
git push origin [tag]
```

### Contributions

While we value open-source contributions to this SDK, this library is generated programmatically.
Feel free to open a PR or a Github issue as a proof of concept and we'll do our best to include it in a future release!
