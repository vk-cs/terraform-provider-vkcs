data "vkcs_compute_flavor" "basic" {
  vcpus = 1
  ram   = 1024
  # specify cpu_generation to distinguish between several flavors with the same CPU and RAM 
  extra_specs = {
    "mcs:cpu_generation" : "cascadelake-v1"
  }
}
