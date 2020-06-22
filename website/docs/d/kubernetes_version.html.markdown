---
layout: "civo"
page_title: "Civo: civo_kubernetes_version"
sidebar_current: "docs-civo-datasource-kubernetes-version"
description: |-
  Get available Civo Kubernetes versions.
---

# civo\_kubernetes\_version

Provides access to the available Civo Kubernetes Service versions, with the ability to filter the results.

## Example Usage

```hcl
data "civo_kubernetes_version" "stable" {
    filter {
        key = "type"
        values = ["stable"]
    }
}
```

### Create a Kubernetes cluster using the most recent version available

```hcl
data "civo_kubernetes_version" "stable" {
    filter {
        key = "type"
        values = ["stable"]
    }
}

resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Traefik"
    num_target_nodes = 4
    kubernetes_version = element(data.civo_kubernetes_version.stable.versions, 0).version
    target_nodes_size = element(data.civo_instances_size.small.sizes, 0).name
}
```

### Pin a Kubernetes cluster to a specific minor version

```hcl
data "civo_kubernetes_version" "minor_version" {
    filter {
        key = "version"
        values = ["0.9.1"]
    }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `version`,
  `label`, `type`, `default`.
* `values` - (Required) Only retrieves the version which keys has value that matches
  one of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `version`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

The following attributes are exported:


* `version` - A version of the kubernetes.
* `label` - The label of this version.
* `type` - The type of the version can be `stable`, `legacy` etc...
* `default` - If is the default version used in all cluster.
