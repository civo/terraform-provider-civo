---
layout: "civo"
page_title: "Civo: civo_disk_image"
sidebar_current: "docs-civo-datasource-disk_image"
description: |-
  Get information on a Civo disk image.
---

# civo\_disk\_image

Get information on an disk image for use in other resources (e.g. creating a Instance).

## Example Usage

This will filter for `debian-10` disk images.

```hcl
data "civo_disk_image" "debian" {
   filter {
        key = "name"
        values = ["debian-10"]
   }
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}
```

It is similar to the previous one with the difference that in this one we use region to only look for the disk images of that region.
The Instance/Kubernetes cluster where you use this data source must be in the same region.

```hcl
data "civo_disk_image" "debian" {
   region = "LON1"
   filter {
        key = "name"
        values = ["debian"]
        match_by = "re"
   }
    sort {
        key = "version"
        direction = "asc"
    }
}

resource "civo_instance" "foo-host" {
    region = "LON1"
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) If is used, them all disk_image will be from that region, has to be declared here if is not declared in the provider
* `filter` - (Optional) Filter the results. The `filter` block is documented below.
* `sort` - (Optional) Sort the results. The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `id`,`name`,`version`,`label`.
* `values` - (Required) Only retrieves the disk_image which keys has value that matches
  one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `id`,`name`,`version`,`label`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.


## Attributes Reference

The following attributes are exported:

* `id` - The id of the disk_image
* `name` - A short human readable name for the disk_image
* `version` - The version of the disk_image.
* `label` - The label of the disk_image.

