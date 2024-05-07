# If the exact amount of disk is not so important to you, then you can specify the minimum value that will satisfy you
# and flavor with minimum of disk will be automatically selected for you.
data "vkcs_compute_flavor" "basic1" {
  vcpus    = 4
  ram      = 8192
  min_disk = 30
}

# But if you also do not specify the exact amount of RAM, then flavor will first be selected based on the minimum RAM
# and if RAM is equal in terms of disk space.
data "vkcs_compute_flavor" "basic2" {
  vcpus    = 4
  min_ram  = 4096
  min_disk = 20
}
