---
layout: "civo"
page_title: "Civo: civo_dns_domain_name"
sidebar_current: "docs-civo-datasource-domain"
description: |-
  Get information on a domain.
---

# civo_dns_domain_name

Get information on a domain. This data source provides the name and the id, this is useful if the domain
name in question is not managed by Terraform.

An error is triggered if the provided domain name is not managed with your
Civo account.

## Example Usage

Get the name and the id file for a domain:

```hcl
data "civo_dns_domain_name" "domain" {
    name = "domain.com"
}

output "domain_output" {
  value = data.civo_dns_domain_name.domain.name
}
output "domain_id_output" {
  value = data.civo_dns_domain_name.domain.id
}

```

```
  $ terraform apply

data.civo_dns_domain_name.domain: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

domain_output = domain.com.
domain_id_output = 6ea98024-c6d7-4d0c-bd01-8ee0cab5224e
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the domain.
* `name` - (Optional) The name of the domain.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a domain.
* `name` - The name of the domain.
