resource "vkcs_baremetal_server" "server" {
  name              = "tf-server-57b4d7f1"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  monitoring        = true
  key_pair          = vkcs_compute_keypair.generated_key.name

  storage_layout {
    disk {
      id   = "system-0"
      type = "SSD"
      size = 447

      partition {
        mount = "/boot/efi"
        fs    = "vfat"
        size  = "512MiB"
      }

      partition {
        mount = "/"
        fs    = "ext4"
        size  = "50GiB"
      }

      partition {
        mount = ""
        fs    = "ext4"
        size  = "60GiB"
      }

      partition {
        fs   = "swap"
        size = "512MiB"
      }
    }

    disk {
      id   = "data-0"
      type = "HDD"
      size = 18500
    }

    disk {
      id   = "data-1"
      type = "HDD"
      size = 18500
    }

    raid {
      id      = "raid-data"
      type    = "raid1"
      members = ["data-0", "data-1"]

      partition {
        mount = "/data"
        fs    = "xfs"
        size  = "10000GiB"
      }
    }
  }

  user_data = <<EOF
    #cloud-config
    package_upgrade: true
    packages:
      - nginx
    runcmd:
      - systemctl start nginx
  EOF

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }
  }

  nic {
    name = "nic1"
    vlan {
      native     = true
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}
