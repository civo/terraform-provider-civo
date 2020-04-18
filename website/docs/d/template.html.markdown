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
you need to utilize any of the image's data.

## Example Usage

```hcl
data "civo_template" "debian" {
    code = "debian-buster"
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = data.civo_size.small.id
    template = data.civo_template.debian.id
}
```
## Argument Reference

The `code` arguments must be provided:

* `code` - A unqiue, alphanumerical, short, human readable code for the template

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

