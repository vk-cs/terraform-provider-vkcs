resource "vkcs_keymanager_secret" "certificate_1" {
  name                 = "certificate"
  payload              = "${file("cert.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "private_key_1" {
  name                 = "private_key"
  payload              = "${file("cert-key.pem")}"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "intermediate_1" {
  name                 = "intermediate"
  payload              = "${file("intermediate-ca.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_container" "tls_1" {
  name = "tls"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${vkcs_keymanager_secret.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${vkcs_keymanager_secret.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${vkcs_keymanager_secret.intermediate_1.secret_ref}"
  }
}

data "vkcs_networking_subnet" "subnet_1" {
  name = "my-subnet"
}

resource "vkcs_lb_loadbalancer" "lb_1" {
  name          = "loadbalancer"
  vip_subnet_id = "${data.vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
  name                      = "https"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 443
  loadbalancer_id           = "${vkcs_lb_loadbalancer.lb_1.id}"
  default_tls_container_ref = "${vkcs_keymanager_container.tls_1.container_ref}"
}
