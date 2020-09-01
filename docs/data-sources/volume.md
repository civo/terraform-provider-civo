---
layout: "civo"
page_title: "Civo: civo_volume"
sidebar_current: "docs-civo-datasource-volume"
description: |-
  Get information on a volume.
---

# civo_volume

Get information on a volume for use in other resources. This data source provides
all of the volumes properties as configured on your Civo account. This is
useful if the volume in question is not managed by Terraform or you need to utilize
any of the volumes data.

An error is triggered if the provided volume name does not exist.

## Example Usage

Get the volume:

```hcl
data "civo_volume" "mysql" {
    name = "database-mysql"
}
```

Reuse the data about a volume to attach it to a Instance:

```hcl
data "civo_volume" "mysql" {
    name = "database-mysql"
}

resource "civo_instance" "mysql-server" {
    hostname = "mysql.domain.com"
    tags = ["mysql", "db"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}

resource "civo_volume_attachment" "foobar" {
    instance_id = civo_instance.mysql-server.id
    volume_id  = data.civo_volume.mysql.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique identifier for the volume.
* `name` - (Optional) The name of the volume.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the volume.
* `name` - Name of the volume.
* `size_gb` - The size of the volume.
* `bootable` - if is bootable or not.
* `mount_point` - The mount point of the volume. 
* `created_at` - The date of the creation of the volume.