data "civo_volume" "myvolume" {
  name = "test-volume-name"
}

output "volume_output" {
  value = data.civo_volume.myvolume
}
