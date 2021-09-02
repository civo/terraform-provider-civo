---
layout: "civo"
page_title: "Civo: civo_ssh_key"
sidebar_current: "docs-civo-resource-ssh-key"
description: |-
  Provides a Civo SSH key resource.
---

# civo\_ssh_key

Provides a Civo SSH Key resource to allow you to manage SSH keys for instance access. Keys created with this resource can be referenced in your instance configuration via their ID.

## Example Usage

```hcl
resource "civo_ssh_key" "my-user"{
    name = "my-user"
    public_key = file("~/.ssh/id_rsa.pub")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSH key for identification
* `public_key` - (Required) The public key. If this is a file, it
can be read using the file interpolation function.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the key
* `name` - The name of the SSH key
* `public_key` - The text of the public key
* `fingerprint` - The fingerprint of the SSH key

## Import

SSH Keys can be imported using the `ssh key id`, e.g.

```
terraform import civo_ssh_key.mykey 87ca2ee4-57d3-4420-b9b6-411b0b4b2a0e
```
