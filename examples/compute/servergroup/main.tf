resource "vkcs_compute_servergroup" "cusom_group" {
  name     = "custom-group-tf-example"
  policies = ["anti-affinity"]
}
