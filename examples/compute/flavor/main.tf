data "vkcs_compute_flavor" "basic" {
  vcpus          = 2
  ram            = 2048
  cpu_generation = "cascadelake-v1"
}