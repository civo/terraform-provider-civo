data "civo_instance" "myhostaname" {
    hostname = "myhostname.com"
}

output "instance_output" {
  value = data.civo_instance.myhostaname.public_ip
}
