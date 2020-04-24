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
        name = "code"
        values = ["buster"]
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
## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `code`,
  `name`.
* `values` - (Required) Only retrieves the template which keys has value that matches
  one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `code`, 
`name`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.


## Attributes Reference

The following attributes are exported:

* `id` - The id of the template
* `code` - A unqiue, alphanumerical, short, human readable code for the template.
* `name` - A short human readable name for the template
* `volume_id` - The ID of a bootable volume, either owned by you or global.
* `image_id` - The Image ID of any default template or the ID of another template.
* `short_description` - A one line description of the template
* `description` - A multi-line description of the template, in Markdown format
* `default_username` - The default username to suggest that the user creates
* `cloud_config` - Commonly referred to as 'user-data', this is a customisation script that is run after
the instance is first booted.

