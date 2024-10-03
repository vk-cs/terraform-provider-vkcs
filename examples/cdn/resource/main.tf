resource "vkcs_cdn_resource" "resource" {
  cname        = local.cname
  origin_group = vkcs_cdn_origin_group.origin_group.id
  options = {
    edge_cache_settings = {
      value = "10m"
    }
    forward_host_header = true
  }
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.certificate.id
  }
}
