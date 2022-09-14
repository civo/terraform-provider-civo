resource "civo_object_store" "backup" {
	name = "backup-server"
	max_size_gb = 500
	region = "LON1"
}
