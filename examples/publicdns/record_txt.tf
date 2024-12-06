resource "vkcs_publicdns_record" "txt" {
  zone_id = vkcs_publicdns_zone.zone.id
  type    = "TXT"
  name    = "@"
  content = "Text example"
  ttl     = 86400
}
