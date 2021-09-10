data "civo_dns_domain_name" "domain" {
    name = "domain.com"
}

data "civo_dns_domain_record" "www" {
    domain_id = data.civo_dns_domain_name.domain.id
    name = "www"
}

output "record_type" {
  value = data.civo_dns_domain_record.www.type
}

output "record_ttl" {
  value = data.civo_dns_domain_record.www.ttl
}
