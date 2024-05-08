# If the exact amount of RAM is not so important to you, then you can specify the minimum value that will satisfy you 
# and flavor with minimum of ram will be automatically selected for you.
data "vkcs_compute_flavor" "basic_min_ram" {
  vcpus   = 4
  min_ram = 4096
}
