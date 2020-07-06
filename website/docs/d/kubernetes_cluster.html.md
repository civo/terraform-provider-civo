---
layout: "civo"
page_title: "Civo: civo_kubernetes_cluster"
sidebar_current: "docs-civo-datasource-kubernetes-cluster"
description: |-
  Get a Civo Kubernetes cluster resource.
---

# civo_kubernetes_cluster

Provides a Civo Kubernetes cluster data source.

**Note:** This data source returns a single kubernetes cluster. When specifying a `name`, an
error is triggered if more than one kubernetes Cluster is found.

## Example Usage

Get the Kubernetes Cluster by name:

```hcl
data "civo_kubernetes_cluster" "my-cluster" {
    name = "my-super-cluster"
}

output "kubernetes_cluster_output" {
  value = data.civo_kubernetes_cluster.my-cluster.master_ip
}
```

Get the Kubernetes Cluster by id:

```hcl
data "civo_kubernetes_cluster" "my-cluster" {
    name = "40ac97ee-b82b-4231-9b60-079c7e2e5d79"
}
```
## Argument Reference

One of following the arguments must be provided:

* `id` - (Optional) The ID of the kubernetes Cluster
* `name` - (Optional) The name of the kubernetes Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Kubernetes cluster.
* `name` - The name of your cluster,.
* `num_target_nodes` - The size of the Kubernetes cluster.
* `target_nodes_size` - The size of each node.
* `kubernetes_version` - The version of Kubernetes.
* `tags` - A list of tags.
* `applications` - A list of application installed.
* `instances` - In addition to the arguments provided, these additional attributes about the cluster's default node instance are exported.
    - `hostname` - The hostname of the instance.
    - `size` - The size of the instance.
    - `region` - The region where instance are.
    - `status` - The status of the instance.
    - `created_at` - The date where the instances was created.
    - `firewall_id` - The firewall id assigned to the instance
    - `public_ip` - The public ip of the instances, only available if the instances is the master
    - `tags` - The tag of the instances
* `installed_applications` - A unique ID that can be used to identify and reference a Kubernetes cluster.
    - `application` - The name of the application
    - `version` - The version of the application
    - `installed` - if installed or not
    - `category` - The category of the application
* `status` - The status of Kubernetes cluster.
* `ready` -If the Kubernetes cluster is ready.
* `kubeconfig` - A representation of the Kubernetes cluster's kubeconfig in yaml format.
* `api_endpoint` - The base URL of the API server on the Kubernetes master node.
* `master_ip` - The Ip of the Kubernetes master node.
* `dns_entry` - The unique dns entry for the cluster in this case point to the master.
* `built_at` - The date where the Kubernetes cluster was build.
* `created_at` - The date where the Kubernetes cluster was create.