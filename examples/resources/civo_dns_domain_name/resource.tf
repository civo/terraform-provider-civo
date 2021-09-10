# Create a new domain name
resource "civo_dns_domain_name" "main" {
  name = "mydomain.com"
}
