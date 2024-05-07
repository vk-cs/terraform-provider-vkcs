# if we set only vcpus = 2 and ram = 2048 we find: STD2-2-2 and STD3-2-2.
# Which differs by process generations.

# If we want a STD2-2-2
data "vkcs_compute_flavor" "basic_cascadelake" {
  vcpus          = 2
  ram            = 2048
  cpu_generation = "cascadelake-v1"
}

# If we want a STD3-2-2
data "vkcs_compute_flavor" "basic_icelake" {
  vcpus          = 2
  ram            = 2048
  cpu_generation = "icelake-v1"
}
