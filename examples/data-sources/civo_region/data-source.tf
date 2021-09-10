data "civo_region" "default" {
    filter {
        key = "default"
        values = ["true"]
    }
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    region = element(data.civo_region.default.regions, 0).code
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}
