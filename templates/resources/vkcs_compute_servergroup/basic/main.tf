resource "vkcs_compute_servergroup" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
