---
layout: "civo"
page_title: "Civo: civo_volume"
sidebar_current: "docs-civo-resource-volume"
description: |-
  Provides a Civo volume resource.
---

# civo\_volume

Provides a Civo Volume which can be attached to an instance in order to provide expanded storage.

## Example Usage

```hcl
# Get network
data "civo_network" "default_network" {
    label = "Default"
}

# Create volume
resource "civo_volume" "db" {
    name = "backup-data"
    size_gb = 5
    network_id = data.civo_network.default_network.id
    depends_on = [
      data.civo_network.default_network
    ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name that you wish to use to refer to this volume.
* `size_gb` - (Required) A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes.
* `network_id` - (Required) The network that the volume belongs to.
* `region` - (Optional) The region for the volume, if not declare we use the region in declared in the provider.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the volume.
* `name` - Name of the volume.
* `size_gb` - The size of the volume.
* `mount_point` - The mount point of the volume.
* `network_id` - The network that the volume belongs to.

## Import

Volumes can be imported using the `volume id`, e.g.

```
terraform import civo_volume.db 506f78a4-e098-11e5-ad9f-000f53306ae1
```
