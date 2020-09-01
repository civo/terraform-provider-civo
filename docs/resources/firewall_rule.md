---
layout: "civo"
page_title: "Civo: civo_firewall_rule"
sidebar_current: "docs-civo-resource-firewall_rule"
description: |-
  Provides a Civo Cloud Firewall Rule resource. This can be used to create, modify, and delete Firewalls Rules.
---

# civo\_firewall_rule

Provides a Civo Cloud Firewall Rule resource. 
This can be used to create, modify, and delete Firewalls Rules.
This resource don't have an update option because the backend don't have the
support for that, so in this case we use ForceNew for all object in the resource.

## Example Usage

```hcl
resource "civo_firewall" "www" {
  name = "www"
}

resource "civo_firewall_rule" "http" {
  firewall_id = civo_firewall.www.id
  protocol = "tcp"
  start_port = "80"
  end_port = "80"
  cidr = ["192.168.1.2/32", "10.10.10.1/32", format("%s/%s",civo_instance.foo.public_ip,"32")]
  direction = "ingress"
  label = "server web"
  depends_on = [civo_firewall.www]
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Required) The Firewall id
* `protocol` (Required) This may be one of "tcp", "udp", or "icmp".
* `start_port` (Required) The start port where traffic to be allowed.
* `end_port` (Required) The end port where traffic to be allowed.
* `cidr` (Required) the IP address of the other end (i.e. not your instance) to affect, or a valid network CIDR (defaults to being globally applied, i.e. 0.0.0.0/0).
* `direction` (Required) will this rule affect ingress traffic
* `label` (Optional) a string that will be the displayed name/reference for this rule (optional)

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Firewall Rule.
* `firewall_id` - The Firewall id
* `protocol` This may be one of "tcp", "udp", or "icmp".
* `start_port` The start port where traffic to be allowed.
* `end_port` The end port where traffic to be allowed.
* `cidr` A list of IP address of the other end (i.e. not your instance) to affect, or a valid network CIDR.
* `direction` Will this rule affect ingress traffic
* `label` A string that will be the displayed name/reference for this rule (optional)

## Import

Firewalls can be imported using the firewall `firewall_id:firewall_rule_id`, e.g.

```
terraform import civo_firewall_rule.http b8ecd2ab-2267-4a5e-8692-cbf1d32583e3:4b0022ee-00b2-4f81-a40d-b4f8728923a7
```
