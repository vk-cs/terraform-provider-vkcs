resource "vkcs_publicdns_record" "a" {
  zone_id = vkcs_publicdns_zone.zone.id
  type = "A"
  name = "test"
  ip = "192.0.2.1"
  ttl = 60
}
