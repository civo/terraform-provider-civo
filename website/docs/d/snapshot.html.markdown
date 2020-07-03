---
layout: "civo"
page_title: "Civo: civo_snapshot"
sidebar_current: "docs-civo-datasource-snapshot"
description: |-
  Get information about a Civo snapshot.
---

# civo_snapshot

Snapshots are saved instances of a block storage volume. Use this data
source to retrieve the ID of a Civo snapshot for use in other
resources.

## Example Usage

Get the snapshot:

```hcl
data "civo_snapshot" "mysql-vm" {
    name = "mysql-vm"
}
```

## Argument Reference

* `id` - (Optional) The ID of the snapshot.
* `name` - (Optional) The name of the snapshot.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Instance snapshot.
* `name` - The name of the snapshot.
* `instance_id` - The ID of the Instance from which the snapshot was be taken.
* `safe` - If is `true` the instance will be shut down during the snapshot if id `false` them not.
* `cron_timing` - A string with the cron format.
* `hostname` - The hostname of the instance.
* `template_id` - The template id.
* `region` - The region where the snapshot was take.
* `size_gb` - The size of the snapshot in GB.
* `state` - The status of the snapshot.
* `next_execution` - if cron was define this date will be the next execution date.
* `requested_at` - The date where the snapshot was requested.
* `completed_at` - The date where the snapshot was completed.