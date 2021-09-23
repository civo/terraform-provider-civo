data "civo_dns_domain_name" "domain" {
    name = "domain.com"
}

output "domain_output" {
  value = data.civo_dns_domain_name.domain.name
}

output "domain_id_output" {
  value = data.civo_dns_domain_name.domain.id
}

