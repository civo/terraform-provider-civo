---
layout: "civo"
page_title: "Civo: civo_instance_size"
sidebar_current: "docs-civo-datasource-instance_size"
description: |-
  Get information on a Civo Instance Size.
---

# civo\_instance\_size

Retrieves information about the Instance sizes that Civo supports,
with the ability to filter the results.

## Example Usage

Most common usage will probably be to supply a size to instances:

```hcl
data "civo_instances_size" "small" {
    filter {
        name = "name"
        values = ["g2.small"]
    }
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = data.civo_template.debian.id
}
```

The data source also supports multiple filters and sorts. For example, to fetch sizes with 1 or 2 virtual CPU and sort by disk:

```hcl
data "civo_instances_size" "small" {
    filter {
        name = "cpu_cores"
        values = [1,2]
    }

    sort {
        name = "disk_gb"
        direction = "desc"
    }

}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = data.civo_template.debian.id
}
```

The data source can also handle multiple sorts. In which case, the sort will be applied in the order it is defined. For example, to sort by memory in ascending order, then sort by disk in descending order between sizes with same memory:

```hcl
data "civo_instances_size" "main" {
  sort {
    // Sort by memory ascendingly
    name       = "ram_mb"
    direction = "asc"
  }

  sort {
    // Then sort by disk descendingly for sizes with same memory
    name       = "disk_gb"
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

* `key` - (Required) Filter the sizes by this key. This may be one of `name`,
  `nice_name`, `cpu_cores`, `ram_mb`, `disk_gb`, `selectable`.
* `values` - (Required) Only retrieves images which keys has value that matches
  one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `name`, 
`nice_name`, `cpu_cores`, `ram_mb`, `disk_gb`, `selectable`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.


## Attributes Reference

The following attributes are exported:

* `name`: The name of the instance size.
* `nice_name`: A human name of the instance size.
* `cpu_cores` - Total of CPU in the instance.
* `ram_mb`: Total of RAM of the instance.
* `disk_gb`: The instance size of SSD.
* `description` - A description of the instance size.
* `selectable`: If can use the instance size.
