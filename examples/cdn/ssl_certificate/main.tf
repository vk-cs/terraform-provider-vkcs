resource "vkcs_cdn_ssl_certificate" "certificate" {
  name        = "tfexample-ssl-certificate"
  certificate = file("${path.module}/certificate.pem")
  private_key = file("${path.module}/private-key.key")
}
