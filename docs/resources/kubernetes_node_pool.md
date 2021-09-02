---
layout: "civo"
page_title: "Civo: civo_kubernetes_node_pool"
sidebar_current: "docs-civo-resource-kubernetes-node-pool"
description: |-
  Provides a Civo Kubernetes cluster node pool resource.
---

# civo\_kubernetes\_node\_pool

Provides a Civo Kubernetes Node Pool resource. While the default node pool must be defined in the `civo_kubernetes_cluster` resource, this resource can be used to add additional ones to a cluster.

## Example Usage

```hcl
# Query xsmall instance size
data "civo_instances_size" "xsmall" {
    filter {
        key = "type"
        values = ["kubernetes"]
    }

    sort {
        key = "ram"
        direction = "asc"
    }
}

# Create a cluster
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    num_target_nodes = 1
    target_nodes_size = element(data.civo_instances_size.xsmall.sizes, 0).name
}

# Add a node pool
resource "civo_kubernetes_node_pool" "front-end" {
   cluster_id = civo_kubernetes_cluster.my-cluster.id
   num_target_nodes = 1
   region = "LON1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kubernetes cluster to which the node pool is associated.
* `region` - (Required) The region of the node pool, has to match that of the cluster.
* `target_nodes_size` - (Optional) The size of each node (The default is currently g3.k3s.medium).
* `num_target_nodes` - (Optional) The number of instances to create (The default at the time of writing is 3).

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

* `id` - A unique ID that can be used to identify and reference a node pool.
* `cluster_id` - (Required) The ID of the Kubernetes cluster to which the node pool is associated.
* `region` - (Required) The region of the node pool, has to match that of the cluster.
* `target_nodes_size` - (Optional) The size of each node.
* `num_target_nodes` - (Optional) The number of instances to create.

## Import

Then the Kubernetes cluster node pool can be imported using the cluster's and pool id `cluster_id:node_pool_id`, e.g.

```
terraform import civo_kubernetes_node_pool.my-pool 1b8b2100-0e9f-4e8f-ad78-9eb578c2a0af:502c1130-cb9b-4a88-b6d2-307bd96d946a
```
