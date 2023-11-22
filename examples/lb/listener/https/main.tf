resource "vkcs_lb_listener" "app_terminated_https" {
  name                      = "app-https-tf-example"
  description               = "Listener for resources/datasources testing"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 8443
  loadbalancer_id           = vkcs_lb_loadbalancer.app.id
  default_tls_container_ref = vkcs_keymanager_container.lb_cert.container_ref
}
