resource "vkcs_keymanager_secret" "certificate" {
  name                 = "certificate-tf-example"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
  payload              = file("${path.module}/certificate.pem")
}
