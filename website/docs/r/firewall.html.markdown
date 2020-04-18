---
layout: "civo"
page_title: "Civo: civo_firewall"
sidebar_current: "docs-civo-resource-firewall"
description: |-
  Provides a Civo Cloud Firewall resource. This can be used to create, modify, and delete Firewalls.
---

# civo\_firewall

Provides a Civo Cloud Firewall resource. This can be used to create,
modify, and delete Firewalls.

## Example Usage

```hcl
resource "civo_firewall" "www" {
  name = "www"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Firewall name


## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Firewall.
* `name` - The name of the Firewall.
* `region` - The region where the firewall was create.

## Import

Firewalls can be imported using the firewall `id`, e.g.

```
terraform import civo_firewall.www b8ecd2ab-2267-4a5e-8692-cbf1d32583e3
```
