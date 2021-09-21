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
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}
