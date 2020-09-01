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

Networks may be looked up by `id` or `label`.

## Example Usage

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

* `id` - The unique identifier of an existing Network.
* `label` - The name of an existing Network.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Network.
* `label` - The label used in the configuration.
* `name` - The name of the network.
* `region` - The region where the network was create.
* `default` - If is the default network.
* `cidr` - The block ip assigned to the network.
