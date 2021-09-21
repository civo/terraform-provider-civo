# Query small instance size
data "civo_instances_size" "small" {
    filter {
        key = "name"
        values = ["g3.small"]
        match_by = "re"
    }

    filter {
        key = "type"
        values = ["instance"]
    }

}

# Query instance disk image
data "civo_disk_image" "debian" {
   filter {
        key = "name"
        values = ["debian-10"]
   }
}

# Create a new instance
resource "civo_instance" "foo" {
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

# Create a new domain name
resource "civo_dns_domain_name" "mydomain" {
  name = "mydomain.com"
}

# Create a new domain record
resource "civo_dns_domain_record" "www" {
    domain_id = civo_dns_domain_name.mydomain.id
    type = "A"
    name = "www"
    value = civo_instance.foo.public_ip
    ttl = 600
    depends_on = [civo_dns_domain_name.mydomain, civo_instance.foo]
}
