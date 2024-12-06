resource "vkcs_publicdns_record" "ns" {
  zone_id = vkcs_publicdns_zone.zone.id
  type    = "NS"
  name    = "@"
  content = "ns1.example.com"
  ttl     = 86400
}
