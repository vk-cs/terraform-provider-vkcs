resource "vkcs_keymanager_container" "lb_cert" {
  name = "container-tf-example"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.priv_key.secret_ref
  }
}
