data "vkcs_keymanager_container" "lb_cert" {
  name = "container-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with container resource example. 
  depends_on = [vkcs_keymanager_container.lb_cert]
}
