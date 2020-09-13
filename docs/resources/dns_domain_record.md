---
layout: "civo"
page_title: "Civo: civo_dns_domain_record"
sidebar_current: "docs-civo-resource-dns-domain-record"
description: |-
  Provides a Civo dns domain record resource.
---

# civo\_dns_domain_record

Provides a Civo dns domain record resource.

## Example Usage

```hcl
# Create a new domain record
resource "civo_dns_domain_record" "www" {
    domain_id = civo_dns_domain_name.main.id
    type = "A"
    name = "www"
    value = civo_instance.foo.public_ip
    ttl = 600
    depends_on = [civo_dns_domain_name.main, civo_instance.foo]  
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required) The id of the domain
* `type` - (Required) The choice of record type from A, CNAME, MX, SRV or TXT
* `name` - (Required) The portion before the domain name (e.g. www) or an @ for the apex/root domain (you cannot use an A record with an amex/root domain)
* `value` - (Required) The IP address (A or MX), hostname (CNAME or MX) or text value (TXT) to serve for this record
* `priority` - (Optional) Useful for MX records only, the priority mail should be attempted it (defaults to 10)
* `ttl` - (Required) How long caching DNS servers should cache this record for, in seconds (the minimum is 600 and the default if unspecified is 600)

## Attributes Reference

The following attributes are exported including the arguments:

* `id` - A unique ID that can be used to identify and reference a Record.
* `account_id` - The id account of the domain
* `created_at` - The date when it was created in UTC format
* `updated_at` - The date when it was updated in UTC format

## Import

Domains can be imported using the `id_domain:id_domain_record`, e.g.

```
terraform import civo_dns_domain_record.www a3cd6832-9577-4017-afd7-17d239fc0bf0:c9a39d14-ee1b-4870-8fb0-a2d4f465e822
```
