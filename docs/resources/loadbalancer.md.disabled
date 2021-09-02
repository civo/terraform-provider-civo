---
layout: "civo"
page_title: "Civo: civo_loadbalancer"
sidebar_current: "docs-civo-resource-loadbalancer"
description: |-
  Provides a Civo Load Balancer resource. This can be used to create, modify, and delete Load Balancers.
---

# civo\_loadbalancer

Provides a Civo Load Balancer resource. This can be used to create,
modify, and delete Load Balancers.

## Example Usage

```hcl
resource "civo_loadbalancer" "myloadbalancer" {
    hostname = "www.foo.com"
    protocol = "http"
    port = 80
    max_request_size = 30
    policy = "round_robin"
    max_conns = 10
    fail_timeout = 40

    backend {
        instance_id = civo_instance.my-test-instance.id
        protocol =  "http"
        port = 80
    }

    backend {
        instance_id = civo_instance.my-test-instance-1.id
        protocol = "http"
        port = "80"
    } 
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The hostname to receive traffic for, e.g. www.example.com (optional: sets hostname to loadbalancer-uuid.civo.com if blank)
* `protocol` - (Required) Either http or https. If you specify https then you must also provide the next two fields, the default is http",
* `tls_certificate` - (Optional) If your protocol is https then you should send the TLS certificate in Base64-encoded PEM format
* `tls_key` - (Optional) If your protocol is https then you should send the TLS private key in Base64-encoded PEM format
* `port` - (Required) You can listen on any port, the default is 80 to match the default protocol of http, if not you must specify it here (commonly 80 for HTTP or 443 for HTTPS)
* `max_request_size` - (Required) The size in megabytes of the maximum request content that will be accepted
* `policy` - (Required) One of: `least_conn` (sends new requests to the least busy server) 
`random` (sends new requests to a random backend), `round_robin` (sends new requests to the next backend in order), 
`ip_hash` (sends requests from a given IP address to the same backend), default is `random`
* `health_check_path` - (Optional) What URL should be used on the backends to determine if it's OK (2xx/3xx status), defaults to /
* `fail_timeout` - (Required) How long to wait in seconds before determining a backend has failed, defaults to 30.
* `max_conns` (Required) - how many concurrent connections can each backend handle, defaults to 10.
* `ignore_invalid_backend_tls` (Optional) - Should self-signed/invalid certificates be ignored from the backend servers, defaults to true.
* `backend` - (Required) A list of backend instances, each containing an `instance_id`, `protocol` (http or https) and `port`.
    - `instance_id` - (Required) - The instance id
    - `protocol` - (Required) - The protocol Either http or https.
    - `port` - (Required) - You can listen on any port.

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

## Import

Load Balancers can be imported using the `id`, e.g.

```
terraform import civo_loadbalancer.myloadbalancer 4de7ac8b-495b-4884-9a69-1050c6793cd6
```
