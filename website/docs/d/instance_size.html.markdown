---
layout: "civo"
page_title: "Civo: civo_instance_size"
sidebar_current: "docs-civo-datasource-instance_size"
description: |-
  Get information on a Civo Instance Size.
---

# civo\_instance\_size

Retrieves information about the Instance sizes that Civo supports,
with the ability to filter the results.

## Example Usage

Get the data about a snapshot:

```hcl
data "civo_instances_size" "disk_200" {
    filter {
        name = "disk"
        values = ["200"]
    }

resource "civo_instance" "example" {
  name   = "example"
  size   = civo_instances_size.disk_200.name
}
```
## Argument Reference

* `filter` - (Required) Filter the results. The filter block is documented below.

`filter` supports the following arguments:

* `name` - - (Required) Filter the sizes by this key. This may be one of `name`, `cpu`, `ram` or `disk`.
* `values` - (Required) Only retrieves images which keys has value that matches one of the values provided here.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the instance size.
* `name`: The name of the instance size.
* `nice_name`: A human name of the instance size.
* `cpu_cores` - Total of CPU in the instance.
* `ram_mb`: Total of RAM of the instance.
* `disk_gb`: The instance size of SSD.
* `description` - A description of the instance size.
* `selectable`: If can use the instance size.
