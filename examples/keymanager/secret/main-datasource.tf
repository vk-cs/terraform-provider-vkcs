data "vkcs_keymanager_secret" "certificate" {
  secret_type = "certificate"
  # This is unnecessary in real life.
  # This is required here to let the example work with secret resource example. 
  depends_on = [vkcs_keymanager_secret.certificate]
}
