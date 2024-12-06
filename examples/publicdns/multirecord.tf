locals {
  google_public_dns_ips = tomap({
    "ip_1" = "8.8.8.8"
    "ip_2" = "8.8.4.4"
  })
}

resource "vkcs_publicdns_record" "multi_a" {
  for_each = local.google_public_dns_ips
  zone_id  = vkcs_publicdns_zone.zone.id
  type     = "A"
  name     = "google-dns-servers"
  ip       = each.value
  ttl      = 60
}
