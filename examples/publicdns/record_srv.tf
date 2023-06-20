resource "vkcs_publicdns_record" "srv" {
  zone_id = vkcs_publicdns_zone.zone.id
  type = "SRV"
  service = "_sip"
  proto = "_udp"
  priority = 10
  weight = 5
  host = "siptarget.com"
  port = 5060
  ttl = 60
}
