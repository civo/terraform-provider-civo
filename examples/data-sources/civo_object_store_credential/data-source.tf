# Read a credential for the object store
data "civo_object_store_credential" "backup" {
	name = "backup-server"
}

# Use the credential to create a bucket
resource "civo_object_store" "backup" {
	name = "backup-server"
	max_size_gb = 500
	region = "LON1"
	access_key_id = data.civo_object_store_credential.backup.access_key_id
}