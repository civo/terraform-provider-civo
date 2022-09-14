# Create a simple credential for the object store
data "civo_object_store_credential" "backup" {
	name = "backup-server"
}

# Create a credential for the object store with a specific access key and secret key
resource "civo_object_store_credential" "backup" {
	name = "backup-server"
	access_key_id = "my-access-key"
	secret_access_key = "my-secret-key"
}

# Use the credential to create a bucket
resource "civo_object_store" "backup" {
	name = "backup-server"
	max_size_gb = 500
	region = "LON1"
	access_key_id = civo_object_store_credential.backup.access_key_id
}