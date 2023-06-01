# Clearblade Terraform Provider

## Project Status

:warning: This project is currently under active development.

Run the following command to build the provider

First, find the GOPATH path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```shell
go env
```

If the GOPATH go environment variable is not set, use the default path, /Users/$Username/go/bin.

Create a new file called .terraformrc in your home directory (~), then add the dev_overrides block below. Change the PATH below to GOPATH/bin/ based on the value returned for GOPATH from the command above.

```shell
provider_installation {

  dev_overrides {
      "clearblade.com/com/clearblade" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Locally install provider and verify with Terraform

```shell
go install .
```

Navigate to the examples directory

```shell
cd examples
```
