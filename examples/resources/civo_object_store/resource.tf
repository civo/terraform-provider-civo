resource "civo_object_store" "backup" {
	name = "backup-server"
	max_size_gb = 500
	region = "LON1"
}

# If you create the bucket without credentials, you can read the credentials in this way
data "civo_object_store_credential" "backup" {
	id = civo_object_store.backup.access_key_id
}