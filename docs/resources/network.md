---
layout: "civo"
page_title: "Civo: civo_network"
sidebar_current: "docs-civo-resource-network"
description: |-
  Provides a Civo Network resource. This can be used to create, modify, and delete Networks.
---

# civo\_loadbalancer

Provides a Civo Network resource. This can be used to create,
modify, and delete Networks.

## Example Usage

```hcl
resource "civo_network" "custom_net" {
    label = "test_network"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The Network label


## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Network.
* `label` - The label used in the configuration.
* `name` - The name of the network.
* `region` - The region where the network was create.
* `default` - If is the default network.
* `cidr` - The block ip assigned to the network.

## Import

Firewalls can be imported using the firewall `id`, e.g.

```
terraform import civo_network.custom_net b8ecd2ab-2267-4a5e-8692-cbf1d32583e3
```