resource "civo_ssh_key" "my-user"{
    name = "my-user"
    public_key = file("~/.ssh/id_rsa.pub")
}
