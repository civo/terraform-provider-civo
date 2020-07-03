---
layout: "civo"
page_title: "Civo: civo_instance"
sidebar_current: "docs-civo-datasource-instance"
description: |-
  Get information on a Instance.
---

# civo_instance

Get information on a Instance for use in other resources. This data source provides
all of the Instance's properties as configured on your Civo account. This
is useful if the Instance in question is not managed by Terraform or you need to
utilize any of the Instance's data.

**Note:** This data source returns a single Instance. When specifying a `hostname`, an
error is triggered if more than one Instance is found.

## Example Usage

Get the Instance by hostname:

```hcl
data "civo_instance" "myhostaname" {
    hostname = "myhostname.com"
}

output "instance_output" {
  value = data.civo_instance.myhostaname.public_ip
}
```

Get the Instance by id:

```hcl
data "civo_instance" "myhostaname" {
    id = "6f283ab7-c37e-42f9-9b4b-f80aea8c006d"
}
```
## Argument Reference

One of following the arguments must be provided:

* `id` - (Optional) The ID of the Instance
* `hostname` - (Optional) The hostname of the Instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Instance.
* `hostname` - The Instance hostname.
* `reverse_dns` - A fully qualified domain name.
* `size` - The name of the size.
* `public_ip_requiered` - This should be either false, true or `move_ip_from:intances_id`.
* `network_id` - This will be the ID of the network.
* `template` - The ID for the template to used to build the instance.
* `initial_user` - The name of the initial user created on the server.
* `notes` - The notes of the instance.
* `sshkey_id` - The ID SSH.
* `firewall_id` - The ID of the firewall used.
* `tags` - An optional list of tags
* `initial_password` - Instance initial password
* `private_ip` - The private ip.
* `public_ip` - The public ip.
* `pseudo_ip` - Is the ip that is used to route the public ip from the internet to the instance using NAT 
* `status` - The status of the instance
* `script` - the contents of a script uploaded
* `created_at` - The date of creation of the instance