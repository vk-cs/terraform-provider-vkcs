resource "vkcs_cdn_resource" "resource" {
  cname        = local.cname # Provide your own value
  origin_group = vkcs_cdn_origin_group.origin_group.id
  options = {
    edge_cache_settings = {
      value = "10m"
    }
    forward_host_header = true
  }
  # Remove if you decided not to enable shielding on the resource
  shielding = {
    enabled = true
    pop_id  = data.vkcs_cdn_shielding_pop.pop.id
  }
  # Remove if not necessary. Specify "lets_encrypt" as the value for
  # the type of the certificate and omit the "id" attribute if you want
  # to issue a Let's Encrypt certificate.
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.certificate.id
  }
}
