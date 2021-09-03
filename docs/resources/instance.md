---
layout: "civo"
page_title: "Civo: civo_instance"
sidebar_current: "docs-civo-resource-instance"
description: |-
  Provides a Civo Instance resource. This can be used to create, modify, and delete Instances.
---

# civo\_instance

Provides a Civo Instance resource. This can be used to create,
modify, and delete Instances.

## Example Usage

```hcl
# Create a new Web instances
resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}
```

```hcl
# Create a new Web instances in a expecific region
resource "civo_instance" "my-test-instance" {
    region = "LON1"
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}
```
## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region for the instance, if not declare we use the region in declared in the provider.
* `hostname` - (Required) The Instance hostname, if is not declare the provider will generate one for you
* `reverse_dns` - (Optional) A fully qualified domain name that should be used as the instance's IP's reverse DNS (optional, uses the hostname if unspecified).
* `size` - (Optional) The name of the size, from the current list, e.g. g3.k3s.small (required).
* `public_ip_required` - (Optional) This should be either `create` or `none` (default: `create`).
* `network_id` - (Optional) This must be the ID of the network from the network listing (optional; default network used when not specified).
* `template` - (Optional) The ID for the template to use to build the instance.
* `initial_user` - (Optional) The name of the initial user created on the server (optional; this will default to the template's default_username and fallback to civo).
* `notes` - (Optional) Add some notes to the instance.
* `sshkey_id` - (Optional) The ID of an already uploaded SSH public key to use for login to the default user (optional; if one isn't provided a random password will be set and returned in the initial_password field).
* `firewall_id` - (Optional) The ID of the firewall to use, from the current list. If left blank or not sent, the default firewall will be used (open to all).
* `script` - (Optional) the contents of a script that will be uploaded to /usr/local/bin/civo-user-init-script on your instance, read/write/executable only by root and then will be executed at the end of the cloud initialization
* `tags` - (Optional) An optional list of tags, represented as a key, value pair.

## Attributes Reference

The following attributes are exported:

* `hostname` - The Instance hostname.
* `reverse_dns` - A fully qualified domain name.
* `size` - The name of the size.
* `cpu_cores` - Total cpu of the inatance.
* `ram_mb` - Total ram of the instance.
* `disk_gb` - The size of the disk.
* `public_ip_requiered` - This should be either `create` or `none`.
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

## Import

Instances can be imported using the instance `id`, e.g.

```
terraform import civo_instance.myintance 18bd98ad-1b6e-4f87-b48f-e690b4fd7413
```
