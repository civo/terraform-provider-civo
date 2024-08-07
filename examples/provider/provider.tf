# Set the variable value in *.tfvars file or using -var="credentials_file=..." CLI flag
variable "credentials_file" {}

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
  credentials_file = var.credentials_file
  region = "LON1"
}

# Create a web server
resource "civo_instance" "web" {
  # ...
}