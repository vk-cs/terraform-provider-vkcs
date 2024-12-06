resource "vkcs_publicdns_record" "aaaa" {
  zone_id = vkcs_publicdns_zone.zone.id
  type    = "AAAA"
  name    = "google-dns-servers"
  ip      = "2001:4860:4860::8888"
  ttl     = 60
}
