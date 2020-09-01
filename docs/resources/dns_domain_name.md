---
layout: "civo"
page_title: "Civo: civo_dns_domain_name"
sidebar_current: "docs-civo-resource-dns-domain-name"
description: |-
  Provides a Civo dns domain name resource.
---

# civo\_dns_domain_name

Provides a Civo dns domain name resource.

## Example Usage

```hcl
# Create a new domain name
resource "civo_dns_domain_name" "main" {
  name = "mydomain.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a domain.
* `name` - The name of the domain.
* `account_id` - The id account of the domain

## Import

Domains can be imported using the `domain name`, e.g.

```
terraform import civo_dns_domain_name.main mydomain.com
```
