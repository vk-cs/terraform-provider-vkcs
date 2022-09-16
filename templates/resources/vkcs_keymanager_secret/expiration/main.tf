resource "vkcs_keymanager_secret" "secret_1" {
  name                 = "certificate"
  payload              = "${file("certificate.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
  expiration           = "${timeadd(timestamp(), format("%dh", 8760))}" # one year in hours

  lifecycle {
    ignore_changes = [
      expiration
    ]
  }
}
