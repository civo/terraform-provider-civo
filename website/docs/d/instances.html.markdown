---
layout: "civo"
page_title: "Civo: instances"
sidebar_current: "docs-civo-datasource-instances"
description: |-
  Retrieve information on Instances.
---

# civo_instances

Get information on Instances for use in other resources, with the ability to filter and sort the results.
If no filters are specified, all Instances will be returned.

This data source is useful if the Instances in question are not managed by Terraform or you need to
utilize any of the Instances' data.

Note: You can use the [`civo_instance`](/docs/providers/civo/d/instance.html) data source to obtain metadata
about a single instance if you already know the `id`, unique `hostname`, or unique `tag` to retrieve.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter images.

For example to find all instances with size `g2.small`:

```hcl
data "civo_instances" "small-size" {
    filter {
        key = "size"
        values = [g2.small]
    }
}
```

You can filter and sort the results as well:

```hcl
data "civo_instances" "small-with-backups" {
  filter {
    key = "size"
    values = [g2.small]
  }
  sort {
    key = "created_at"
    direction = "desc"
  }
}
```
if you don't know the size you can use the [`civo_instances_size`](/docs/providers/civo/d/instances_size.html) data source to obtain metadata
and use in this way:

```hcl
data "civo_instances_size" "small" {
    filter {
        key = "name"
        values = ["small"]
    }
}

data "civo_instances" "small-with-backups" {
  filter {
    key = "size"
    values = [data.civo_instances_size.small.sizes[1].name]
  }
  sort {
    key = "created_at"
    direction = "desc"
  }
}
```


## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the Instances by this key. This may be one of '`id`, `hostname`, `public_ip`, `private_ip`,
  `pseudo_ip`, `size`, `template` or `created_at`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves Instances
  where the `key` field takes on one or more of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the Instance by this key. This may be one of `id`, `hostname`, `public_ip`, `private_ip`,
  `pseudo_ip`, `size`, `template` or `created_at`.

* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `instances` - A list of Instances satisfying any `filter` and `sort` criteria. Each instance has the following attributes:  

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