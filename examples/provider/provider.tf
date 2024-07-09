# Set the variable value in *.tfvars file or using -var="credential_file=..." CLI flag
variable "credential_file" {}

# Specify required provider as maintained by civo
terraform {
  required_providers {
    civo = {
      source = "civo/civo"
    }
  }
}

# Configure the Civo Provider
provider "civo" {
  credential_file = var.credential_file
  region = "LON1"
}

# Create a web server
resource "civo_instance" "web" {
  # ...
}