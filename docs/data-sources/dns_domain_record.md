---
layout: "civo"
page_title: "Civo: civo_dns_domain_record"
sidebar_current: "docs-civo-datasource-record"
description: |-
  Get information on a DNS record.
---

# civo_dns_domain_record

Get information on a DNS record. This data source provides the name, TTL, and zone
file as configured on your Civo account. This is useful if the record
in question is not managed by Terraform.

An error is triggered if the provided domain name or record are not managed with
your Civo account.

## Example Usage

Get data from a DNS record:

```hcl
data "civo_dns_domain_name" "domain" {
    name = "domain.com"
}

data "civo_dns_domain_record" "www" {
    domain_id = data.civo_dns_domain_name.domain.id
    name = "www"
}

output "record_type" {
  value = data.civo_dns_domain_record.www.type
}

output "record_ttl" {
  value = data.civo_dns_domain_record.www.ttl
}
```

```
  $ terraform apply

data.civo_dns_domain_record.www: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

record_ttl = 3600
record_type = A
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the record.
* `domain_id` - (Required) The domain id of the record.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Record.
* `domain_id` - The id of the domain
* `type` - The choice of record type from A, CNAME, MX, SRV or TXT
* `name` - The portion before the domain name (e.g. www) or an @ for the apex/root domain (you cannot use an A record with an amex/root domain)
* `value` - The IP address (A or MX), hostname (CNAME or MX) or text value (TXT) to serve for this record
* `priority` - The priority of the record.
* `ttl` - How long caching DNS servers should cache this record.
* `account_id` - The id account of the domain.
* `created_at` - The date when it was created in UTC format
* `updated_at` - The date when it was updated in UTC format