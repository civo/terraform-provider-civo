# Set the variable value in *.tfvars file or using -var="civo_token=..." CLI flag
variable "civo_token" {}

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
  token = var.civo_token
  region = "LON1"
}

# Create a web server
resource "civo_instance" "web" {
  # ...
}
