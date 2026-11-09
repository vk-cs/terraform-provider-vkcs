resource "vkcs_compute_instance" "restored" {
  count      = var.enable_restore ? 1 : 0
  name       = var.instance_name
  restore_id = var.restore_id
}
