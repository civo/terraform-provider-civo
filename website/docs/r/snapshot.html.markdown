---
layout: "civo"
page_title: "Civo: civo_snapshot"
sidebar_current: "docs-civo-resource-snapshot"
description: |-
  Provides a Civo Instance snapshot resource.
---

# civo\_snapshot

Provides a resource which can be used to create a snapshot from an existing Civo Instance.

## Example Usage

```hcl
resource "civo_snapshot" "myinstance-backup" {
    name = "myinstance-backup"
    instance_id = civo_instance.myinstance.id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the instance snapshot.
* `instance_id` - (Required) The ID of the Instance from which the snapshot will be taken.
* `safe` - (Optional) If `true` the instance will be shut down during the snapshot to ensure all files 
are in a consistent state (e.g. database tables aren't in the middle of being optimised 
and hence risking corruption). The default is `false` so you experience no interruption 
of service, but a small risk of corruption.
* `cron_timing` - (Optional) If a valid cron string is passed, the snapshot will be saved as an automated snapshot 
continuing to automatically update based on the schedule of the cron sequence provided 
The default is nil meaning the snapshot will be saved as a one-off snapshot.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Droplet snapshot.
* `name` - The name of the snapshot.
* `instance_id` - The ID of the Instance from which the snapshot was be taken.
* `safe` - If is `true` the instance will be shut down during the snapshot if id `false` them not.
* `cron_timing` - A string with the cron format.
* `hostname` - The hostname of the instance.
* `template_id` - The template id.
* `region` - The region where the snapshot was take.
* `size_gb` - The size of the snapshot in GB.
* `state` - The status of the snapshot.
* `next_execution` - if cron was define the this date will be the next execution date.
* `requested_at` - The date where the snapshot was requested.
* `completed_at` - The date where the snapshot was completed.


## Import

Instance Snapshots can be imported using the `snapshot id`, e.g.

```
terraform import civo_snapshot.myinstance-backup 4cc87851-e1d0-4270-822a-b36d28c7a77f
```
