resource "vkcs_publicdns_zone" "zone" {
  zone = local.zone_name
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
}
