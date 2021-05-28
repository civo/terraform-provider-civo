---
layout: "civo"
page_title: "Civo: civo_network"
sidebar_current: "docs-civo-datasource-network"
description: |-
  Get information about a Network.
---

# civo_network

Retrieve information about a Network for use in other resources.

This data source provides all of the Network's properties as configured on your
Civo account. This is useful if the Network in question is not managed by
Terraform or you need to utilize any of the Network's data.

Networks may be looked up by `id` or `label`, and optional you can pass `region`
if you wanna made a lookup for an expecific network inside that region.
## Example Usage

### Network by name in a region

```hcl
data "civo_network" "test" {
    label = "test-network"
    region = "NYC1"
}
```

### Network By Name

```hcl
data "civo_network" "test" {
    label = "test-network"
}
```

Reuse the data about a Network to assign a Instance to it:

```hcl
data "civo_network" "test" {
    label = "test-network"
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = data.civo_size.small.id
    template = data.civo_template.debian.id
    network_id = data.civo_network.test.id
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - (Optional) The unique identifier of an existing Network.
* `label` - (Optional) The label of an existing Network.
* `region` - (Optional) The region of an existing Network.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Network.
* `label` - The label used in the configuration.
* `name` - The name of the network.
* `default` - If is the default network.
