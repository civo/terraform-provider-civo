---
layout: "civo"
page_title: "Civo: civo_region"
sidebar_current: "docs-civo-datasource-region"
description: |-
  Get information on a Civo Region.
---

# civo_region

Retrieves information about the Region that Civo supports,
with the ability to filter the results.

## Example Usage

Most common usage will probably be to supply regions:

```hcl
data "civo_region" "default" {
    filter {
        key = "default"
        values = ["true"]
    }
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    region = element(data.civo_region.default.regions, 0).code
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}
```

The data source also supports multiple filters and sorts. For example, to fetch all region for `us` and `uk`:

```hcl
data "civo_region" "NYC" {
    filter {
        key = "code"
        values = ["NYC1"]
    }

    sort {
        key = "code"
        direction = "desc"
    }

}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `code`, `name`, `country`, `default`.
* `values` - (Required) Only retrieves region which keys has value that matches one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `code`,`name`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

The following attributes are exported:

* `code`: The code of the region.
* `name`: A human name of the region.
* `country` The country of the region.
* `default`: If the region is the default region.
