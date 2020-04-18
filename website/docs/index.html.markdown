---
layout: "civo"
page_title: "Provider: Civo"
sidebar_current: "docs-civo-index"
description: |-
  The Civo provider is used to interact with the resources supported by Civo. The provider needs to be configured with the proper credentials before it can be used.
---

# Civo Provider

The Civo provider is used to interact with the
resources supported by Civo. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Set the variable value in *.tfvars file
# or using -var="civo_token=..." CLI option
variable "civo_token" {}

# Configure the Civo Provider
provider "civo" {
  token = var.civo_token
}

# Create a web server
resource "civo_instance" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) This is the Civo API token. Alternatively, this can also be specified
  using environment variables ordered by precedence:
  * `CIVO_TOKEN`