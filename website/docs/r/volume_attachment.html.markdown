---
layout: "civo"
page_title: "Civo: civo_volume_attachment"
sidebar_current: "docs-civo-resource-volume-attachment"
description: |-
  Provides a Civo volume attachment resource.
---

# civo\_volume\_attachment

Manages attaching a Volume to a Instance.

## Example Usage

```hcl
resource "civo_volume" "db" {
     name = "backup-data"
     size_gb = 60
     bootable = false
}

resource "civo_volume_attachment" "foobar" {
  instance_id = civo_instance.my-test-instance.id
  volume_id  = civo_volume.db.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the instance to attach the volume to.
* `volume_id` - (Required) ID of the Volume to be attached to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the volume attachment.