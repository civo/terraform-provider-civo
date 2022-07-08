# Send to create a reserved IP
resource "civo_reserved_ip" "www" {
    name = "nginx-www" 
}

# We assign the reserved IP to the instance
resource "civo_instance_reserved_ip_assignment" "webserver-www" {
  instance_id = civo_instance.www.id
  reserved_ip_id = civo_reserved_ip.web-server.id
}