---
layout: "civo"
page_title: "Civo: civo_ssh_key"
sidebar_current: "docs-civo-datasource-ssh-key"
description: |-
  Get information on a ssh key.
---

# civo_ssh_key

Get information on a ssh key. This data source provides the name,
and fingerprint as configured on your Civo account. This is useful if
the ssh key in question is not managed by Terraform or you need to utilize any
of the keys data.

An error is triggered if the provided ssh key name does not exist.

## Example Usage

Get the ssh key:

```hcl
data "civo_ssh_key" "example" {
  name = "example"
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
    sshkey_id = data.civo_ssh_key.example.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the ssh key.
* `name` - (Optional) The name of the ssh key.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the ssh key.
* `name`: The name of the ssh key.
* `fingerprint`: The fingerprint of the public key of the ssh key.
