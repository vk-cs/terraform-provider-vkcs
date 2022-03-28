terraform {
  required_providers {
    vkcs = {
      source  = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}

resource "vkcs_blockstorage_volume" "bs-volume" {
  name = "bs-volume"
  size = 8
  volume_type = "ceph-hdd"
  availability_zone = "DP1"
}