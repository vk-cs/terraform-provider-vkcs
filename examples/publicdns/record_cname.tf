resource "vkcs_publicdns_record" "cname" {
  zone_id = vkcs_publicdns_zone.zone.id
  type    = "CNAME"
  name    = "example"
  content = "example.com"
  ttl     = 86400
}
