---
layout: "civo"
page_title: "Civo: civo_kubernetes_version"
sidebar_current: "docs-civo-datasource-kubernetes-version"
description: |-
  Get available Civo Kubernetes versions.
---

# civo\_kubernetes\_version

Provides access to the available Civo Kubernetes Service versions.

## Example Usage

```hcl
data "civo_kubernetes_version" "stable" {
    filter {
        name = "type"
        values = ["stable"]
    }
}
```

### Create a Kubernetes cluster using the most recent version available

```hcl
data "civo_kubernetes_version" "stable" {
    filter {
        name = "type"
        values = ["stable"]
    }
}

resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Traefik"
    num_target_nodes = 4
    kubernetes_version = data.civo_kubernetes_version.stable.id
    target_nodes_size = data.civo_size.small.id
}
```

### Pin a Kubernetes cluster to a specific minor version

```hcl
data "civo_kubernetes_version" "minor_version" {
    filter {
        name = "version"
        values = ["0.9.1"]
    }
}
```

## Argument Reference

* `filter` - (Required) Filter the results. The filter block is documented below.

`filter` supports the following arguments:

* `name` - - (Required) Filter the sizes by this key. This may be one of `version`, `type`.
* `values` - (Required) Only retrieves images which keys has value that matches one of the values provided here.

## Attributes Reference

The following attributes are exported:

* `id` - The id used when you create a new cluster is the same like `version`
* `version` - A version of the kubernetes.
* `label` - The label of this version.
* `type` - The type of the version can be `stable`, `legacy` etc...
* `default` - If is the default version used in all cluster.
