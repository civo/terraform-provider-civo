---
layout: "civo"
page_title: "Civo: civo_kubernetes_cluster"
sidebar_current: "docs-civo-resource-kubernetes-cluster"
description: |-
  Provides a Civo Kubernetes cluster resource.
---

# civo\_kubernetes\_cluster

Provides a Civo Kubernetes cluster resource. This can be used to create, delete, and modify clusters.

## Example Usage

```hcl
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Portainer, Traefik"
    num_target_nodes = 4
    kubernetes_version = element(data.civo_kubernetes_version.stable.versions, 0).version
    target_nodes_size = element(data.civo_instances_size.small.sizes, 0).name
}
```

### Kubernetes Terraform Provider Example

The cluster's kubeconfig is exported as an attribute allowing you to use it with the [Kubernetes Terraform provider](https://www.terraform.io/docs/providers/kubernetes/index.html). For example:

```hcl
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Portainer, Traefik"
    num_target_nodes = 4
    kubernetes_version = element(data.civo_kubernetes_version.stable.versions, 0).version
    target_nodes_size = element(data.civo_instances_size.small.sizes, 0).name
}

provider "kubernetes" {
  load_config_file = false
  host  = civo_kubernetes_cluster.my-cluster.api_endpoint
  username = yamldecode(civo_kubernetes_cluster.my-cluster.kubeconfig).users[0].user.username
  password = yamldecode(civo_kubernetes_cluster.my-cluster.kubeconfig).users[0].user.password
  cluster_ca_certificate = base64decode(
    yamldecode(civo_kubernetes_cluster.my-cluster.kubeconfig).clusters[0].cluster.certificate-authority-data
  )
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Kubernetes cluster.
* `num_target_nodes` - (Optional) The number of instances to create (optional, the default at the time of writing is 3).
* `target_nodes_size` - (Optional) The size of each node (optional, the default is currently g2.small)
* `kubernetes_version` - (Optional) The version of k3s to install (optional, the default is currently the latest available).
* `tags` - (Optional) A space separated list of tags, to be used freely as required.
* `applications` - (Optional) A comma separated list of applications to install. Spaces within application names are fine, but shouldn't be either side of the comma. If you want to remove a default installed application, prefix it with a '-', e.g. -traefik

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

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
* `dns_entry` - The unique dns entry for the cluster in this case point to the master.
* `built_at` - The date where the Kubernetes cluster was build.
* `created_at` - The date where the Kubernetes cluster was create.


## Import

Then the Kubernetes cluster can be imported using the cluster's `id`, e.g.

```
terraform import civo_kubernetes_cluster.my-cluster 1b8b2100-0e9f-4e8f-ad78-9eb578c2a0af
```
