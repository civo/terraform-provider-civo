---
layout: "civo"
page_title: "Civo: civo_firewall_rule"
sidebar_current: "docs-civo-resource-firewall_rule"
description: |-
  Provides a Civo Cloud Firewall Rule resource. This can be used to create, modify, and delete Firewalls Rules.
---

# civo\_firewall_rule

Provides a Civo Firewall Rule resource. This can be used to create, modify, and delete firewalls rules. This resource don't have an update option because the backend don't have the support for that, so in this case we use ForceNew for all object in the resource.

## Example Usage

```hcl
# Create a firewall
resource "civo_firewall" "www" {
    name = "www"
}

# Create a firewall rule
resource "civo_firewall_rule" "http" {
    firewall_id = civo_firewall.www.id
    protocol = "tcp"
    start_port = "80"
    end_port = "80"
    cidr = ["192.168.1.2/32"]
    direction = "ingress"
    label = "web-server"
    depends_on = [civo_firewall.www]
}
```

## Example Usage With Other Resources

```hcl
# Query small instance size
data "civo_instances_size" "small" {
    filter {
        key = "name"
        values = ["g3.small"]
        match_by = "re"
    }

    filter {
        key = "type"
        values = ["instance"]
    }

}

# Query instance template
data "civo_template" "debian" {
   filter {
        key = "name"
        values = ["debian-10"]
   }
}

# Create a new instance
resource "civo_instance" "foo" {
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}

# Create a network
resource "civo_network" "custom_net" {
    label = "my-custom-network"
}

# Create a firewall
resource "civo_firewall" "custom_firewall" {
  name = "my-custom-firewall"
  network_id = civo_network.custom_net.id
}

# Create a firewall rule and only allow
# connections from instance we created above
resource "civo_firewall_rule" "custom_port" {
    firewall_id = civo_firewall.custom_firewall.id
    protocol = "tcp"
    start_port = "3000"
    end_port = "3000"
    cidr = [format("%s/%s",civo_instance.foo.public_ip,"32")]
    direction = "ingress"
    label = "custom-application"
    depends_on = [civo_firewall.custom_firewall]
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Required) The Firewall id
* `protocol` - (Required) This may be one of "tcp", "udp", or "icmp".
* `start_port` - (Required) The start port where traffic to be allowed.
* `end_port` - (Required) The end port where traffic to be allowed.
* `cidr` - (Required) The CIDR notation of the other end to affect, or a valid network CIDR (e.g. 0.0.0.0/0 to open for everyone or 1.2.3.4/32 to open just for a specific IP address.
* `direction` - (Required) Will this rule affect ingress traffic
* `label` - (Optional) A string that will be the displayed name/reference for this rule
* `region` - (Optional) Region for the rule, if is not defined we use the global defined in the provider

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Firewall Rule.
* `firewall_id` - The Firewall id
* `protocol` - This may be one of "tcp", "udp", or "icmp".
* `start_port` - The start port where traffic to be allowed.
* `end_port` - The end port where traffic to be allowed.
* `cidr` - A list of CIDR notations of the other end to affect.
* `direction` - Will this rule affect ingress traffic
* `label` - A string that will be the displayed name/reference for this rule

## Import

Firewalls can be imported using the firewall `firewall_id:firewall_rule_id`, e.g.

```
terraform import civo_firewall_rule.http b8ecd2ab-2267-4a5e-8692-cbf1d32583e3:4b0022ee-00b2-4f81-a40d-b4f8728923a7
```
