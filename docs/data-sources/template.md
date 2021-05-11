---
layout: "civo"
page_title: "Civo: civo_template"
sidebar_current: "docs-civo-datasource-template"
description: |-
  Get information on a Civo template.
---

# civo\_template

Get information on an template for use in other resources (e.g. creating a Instance).
This is useful if the template in question is not managed by Terraform or 
you need to utilize any of the image's data, with the ability to filter the results.

## Example Usage

```hcl
data "civo_template" "debian" {
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
    template = element(data.civo_template.debian.templates, 0).id
}
```

This filter will garatice to install the latest version of debian always

```hcl
data "civo_template" "debian" {
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
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) If is used, them all instances will be from that region.
* `filter` - (Optional) Filter the results. The `filter` block is documented below.
* `sort` - (Optional) Sort the results. The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `id`,`name`,`version`,`label`.
* `values` - (Required) Only retrieves the template which keys has value that matches
  one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `id`,`name`,`version`,`label`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.


## Attributes Reference

The following attributes are exported:

* `id` - The id of the template
* `name` - A short human readable name for the template
* `version` - The version of the template.
* `label` - The label of the template.

