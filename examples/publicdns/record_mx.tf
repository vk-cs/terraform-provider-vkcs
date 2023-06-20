resource "vkcs_publicdns_record" "mx" {
  zone_id = vkcs_publicdns_zone.zone.id
  type = "MX"
  name = "@"
  priority = 10
  content = "mx.example.com"
  ttl = 86400
}
