resource "vkcs_networking_floatingip" "base_fip" {
  pool        = "ext-net"
  description = "floating ip in external net tf example"
}
