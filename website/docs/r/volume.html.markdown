---
layout: "civo"
page_title: "Civo: civo_volume"
sidebar_current: "docs-civo-resource-volume"
description: |-
  Provides a Civo volume resource.
---

# civo\_volume

Provides a Civo volume which can be attached to a Instance in order to provide expanded storage.

## Example Usage

```hcl
resource "civo_volume" "db" {
     name = "backup-data"
     size_gb = 60
     bootable = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name that you wish to use to refer to this volume .
* `size_gb` - (Required) A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes .
* `bootable` - (Required) Mark the volume as bootable.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the volume.
* `name` - Name of the volume.
* `size_gb` - The size of the volume.
* `bootable` - if is bootable or not.
* `mount_point` - The mount point of the volume. 
* `created_at` - The date of the creation of the volume.

## Import

Volumes can be imported using the `volume id`, e.g.

```
terraform import civo_volume.db 506f78a4-e098-11e5-ad9f-000f53306ae1
```
