resource "vkcs_networking_floatingip" "base_fip" {
  pool        = "internet"
  description = "floating ip in external net tf example"
}
