---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "civo_reserved_ip Resource - terraform-provider-civo"
subcategory: "Civo Network"
description: |-
  Provides a Civo reserved IP to represent a publicly-accessible static IP addresses that can be mapped to one of your Instancesor Load Balancer.
---

# civo_reserved_ip (Resource)

Provides a Civo reserved IP to represent a publicly-accessible static IP addresses that can be mapped to one of your Instancesor Load Balancer.

## Example Usage

```terraform
resource "civo_reserved_ip" "www" {
    name = "nginx-www" 
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name for the ip address

### Optional

- `region` (String) The region of the ip

### Read-Only

- `id` (String) The ID of this resource.
- `ip` (String) The IP Address of the resource

## Import

Import is supported using the following syntax:

```shell
terrafom import civo_reserved_ip.www 9f0e86fc-b2c6-46b4-82ed-2f28419f8ae3
```
