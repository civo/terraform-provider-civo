data "civo_ssh_key" "example" {
  name = "example"
}

resource "civo_instance" "my-test-instance" {
    hostname = "foo.com"
    tags = ["python", "nginx"]
    notes = "this is a note for the server"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
    sshkey_id = data.civo_ssh_key.example.id
}
