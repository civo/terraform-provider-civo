---
layout: "civo"
page_title: "Civo: civo_loadbalancer"
sidebar_current: "docs-civo-datasource-loadbalancer"
description: |-
  Get information on a loadbalancer.
---

# civo_loadbalancer

Get information on a load balancer for use in other resources. This data source
provides all of the load balancers properties as configured on your Civo
account. This is useful if the load balancer in question is not managed by
Terraform or you need to utilize any of the load balancers data.

An error is triggered if the provided load balancer name does not exist.

## Example Usage

Get the load balancer:

```hcl
data "civo_loadbalancer" "example" {
    hostname = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional)  The ID of the Load Balancer.
* `hostname` - (Optional) The hostname of the Load Balancer.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Load Balancer
* `hostname` - The hostname of the Load Balancer
* `protocol` - The protocol used
* `tls_certificate` - If is set will be returned
* `tls_key` - If is set will be returned
* `port` - The port set in the configuration
* `max_request_size` - The max request size set in the configuration
* `policy` - The policy set in the Load Balancer
* `health_check_path` - The path to check the health of the backend
* `fail_timeout` - The wait time until the backend is marked as a failure
* `max_conns` - How many concurrent connections can each backend handle
* `ignore_invalid_backend_tls` - Should self-signed/invalid certificates be ignored from the backend servers
* `backend` - A list of backend instances
     - `instance_id` - The instance id
     - `protocol` - The protocol used in the configuration.
     - `port` - The port set in the configuration.