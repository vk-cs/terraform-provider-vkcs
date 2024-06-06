resource "vkcs_dc_ip_port_forwarding" "dc-ip-port-forwarding" {
  dc_interface_id = vkcs_dc_interface.dc_interface.id
  name            = "tf-example"
  description     = "tf-example-description"
  protocol        = "udp"
  to_destination  = "172.17.20.30"
}
