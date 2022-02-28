# TODO
data civo_loadbalancer "my-lb" {
  #id = "c385638f-6bb7-4d74-840c-4d98f3d15082" // Optional
  name = "lb-name"
}

output "civo_loadbalancer_output" {
  value = data.civo_loadbalancer.my-lb.public_ip
}