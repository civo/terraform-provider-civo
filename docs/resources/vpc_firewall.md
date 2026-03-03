---
page_title: "civo_vpc_firewall Resource - terraform-provider-civo"
subcategory: "Civo VPC"
description: |-
  Provides a Civo VPC firewall resource. This can be used to create, modify, and delete firewalls.
---

# civo_vpc_firewall (Resource)

Provides a Civo VPC firewall resource. This can be used to create, modify, and delete firewalls.

~> **Note:** This resource replaces the deprecated `civo_firewall` resource and uses the VPC-scoped API endpoints.

## Example Usage

### Custom ingress and egress rules firewall

```terraform
resource "civo_vpc_network" "example" {
  label = "example-network"
}

resource "civo_vpc_firewall" "example" {
  name                 = "example-firewall"
  network_id           = civo_vpc_network.example.id
  create_default_rules = false
  ingress_rule {
    label      = "http"
    protocol   = "tcp"
    port_range = "80"
    cidr       = ["0.0.0.0/0"]
    action     = "allow"
  }

  egress_rule {
    label      = "all"
    protocol   = "tcp"
    port_range = "1-65535"
    cidr       = ["0.0.0.0/0"]
    action     = "allow"
  }
}
```

### Simple firewall

```terraform
resource "civo_vpc_firewall" "example" {
    name       = "example-firewall"
    network_id = civo_vpc_network.example.id
}
```

## Argument Reference

### Required

- `name` (String) The firewall name

### Optional

- `create_default_rules` (Boolean) The create rules flag is used to create the default firewall rules, if is not defined will be set to true, and if you set to false you need to define at least one ingress or egress rule.
- `egress_rule` (Block Set) The egress rules, this is a list of rules that will be applied to the firewall (see [below for nested schema](#nestedblock--egress_rule))
- `ingress_rule` (Block Set) The ingress rules, this is a list of rules that will be applied to the firewall (see [below for nested schema](#nestedblock--ingress_rule))
- `network_id` (String) The firewall network, if is not defined we use the default network
- `region` (String) The firewall region, if is not defined we use the global defined in the provider

<a id="nestedblock--ingress_rule"></a>
### Nested Schema for `ingress_rule`

Required:

- `action` (String) The action of the rule can be allow or deny.
- `cidr` (Set of String) The CIDR notation of the other end to affect.

Optional:

- `label` (String) A string that will be the displayed name/reference for this rule
- `port_range` (String) The port or port range to open
- `protocol` (String) The protocol choice from `tcp`, `udp` or `icmp` (the default if unspecified is `tcp`)

Read-Only:

- `id` (String) The ID of the firewall rule.

<a id="nestedblock--egress_rule"></a>
### Nested Schema for `egress_rule`

Required:

- `action` (String) The action of the rule can be allow or deny.
- `cidr` (Set of String) The CIDR notation of the other end to affect.

Optional:

- `label` (String) A string that will be the displayed name/reference for this rule
- `port_range` (String) The port or port range to open
- `protocol` (String) The protocol choice from `tcp`, `udp` or `icmp` (the default if unspecified is `tcp`)

Read-Only:

- `id` (String) The ID of the firewall rule.

## Import

```shell
terraform import civo_vpc_firewall.www b8ecd2ab-2267-4a5e-8692-cbf1d32583e3
```
