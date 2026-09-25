---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_server"
description: |-
  Manages a bare metal server resource within VKCS.
---

# vkcs_baremetal_server



## Example Usage
### Basic Server
```terraform
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
```

### Server with bonded interfaces
```terraform
resource "vkcs_baremetal_server" "server_bond" {
  name              = "server-bond"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  monitoring        = true
  key_pair          = vkcs_compute_keypair.generated_key.name

  bond {
    name            = "bond0"
    interface_names = ["nic0", "nic1"]
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }

    vlan {
      id         = 100
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}
```

```terraform
resource "vkcs_baremetal_server" "server_bond" {
  name              = "server-bond"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  monitoring        = true
  key_pair          = vkcs_compute_keypair.generated_key.name

  bond {
    name            = "bond0"
    interface_names = ["nic0"]
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }
  }

  bond {
    name            = "bond1"
    interface_names = ["nic1"]

    vlan {
      native     = true
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}
```

### Server with tagged VLAN interfaces
```terraform
resource "vkcs_baremetal_server" "server_vlan" {
  name              = "server-vlan"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  monitoring        = true
  key_pair          = vkcs_compute_keypair.generated_key.name

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }

    vlan {
      id         = 100
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}
```

## Argument Reference
- `flavor_id` **required** *string* &rarr;  Server flavor to rent.

- `key_pair` **required** *string* &rarr;  The name of a key pair to put on the server. The key pair must already be created and associated with the tenant's account. Changing this creates a new server.

- `name` **required** *string* &rarr;  Name of the bare metal server.

- `availability_zone` optional *string* &rarr;  Availability zone. If not specified, we will chose the availability zone for you.

- `bond` optional &rarr;  Link aggregation interfaces (bonds).
    - `interface_names` **required** *string* &rarr;  List of interface names participating in the bond.

    - `name` **required** *string* &rarr;  Bond interface name (e.g. bond0).

    - `vlan` optional &rarr;  VLAN configuration applied to the bond.
        - `id` optional *number* &rarr;  Number of the VLAN.

        - `native` optional *boolean* &rarr;  Whether the VLAN is native.

        - `network_id` optional *string* &rarr;  ID of the network.

        - `subnet_id` optional *string* &rarr;  ID of the subnet.

- `monitoring` optional *boolean* &rarr;  Whether the monitoring is actively enabled.

- `nic` optional &rarr;  Physical network interfaces.
    - `name` **required** *string* &rarr;  Interface name (e.g. nic0, eno1). Acts as unique identifier.

    - `vlan` optional &rarr;  VLAN configuration. Allowed only if interface is not part of a bond.
        - `id` optional *number* &rarr;  Number of the VLAN.

        - `native` optional *boolean* &rarr;  Whether the VLAN is native.

        - `network_id` optional *string* &rarr;  ID of the network.

        - `subnet_id` optional *string* &rarr;  ID of the subnet.

- `os_id` optional *string* &rarr;  Set os id.

- `region` optional *string* &rarr;  The region to fetch the bare metal server from, defaults to the provider's region.

- `storage_layout` optional &rarr;  Storage layout of the bare metal server: disks carry their own partitions, raids are assembled from whole disks. Changing this triggers reprovisioning.
    - `disk` optional &rarr;  Logical disks and their partition layout.
        - `id` **required** *string* &rarr;  Logical disk identifier.

        - `size` **required** *number* &rarr;  Declared disk size in whole GiB, taken from the flavor. Used to pick a real disk within the size tolerance.

        - `type` **required** *string* &rarr;  Storage medium of the disk: SSD, HDD or NVME (case-insensitive). Must match the flavor disk type.

        - `partition` optional &rarr;  Ordered partitions of the device; order determines placement on disk.
            - `fs` **required** *string* &rarr;  Filesystem type: ext4, xfs, vfat or swap (case-insensitive). Swap requires an empty mount.

            - `size` **required** *string* &rarr;  Partition size with an IEC suffix, for example 512MiB or 50GiB.

            - `mount` optional *string* &rarr;  Mount point of the partition. Empty or omitted means the partition is created but not mounted.

    - `raid` optional &rarr;  RAID arrays assembled from whole disks; partitions are cut on top of the md device.
        - `id` **required** *string* &rarr;  RAID identifier.

        - `members` **required** *string* &rarr;  Disk identifiers the RAID is assembled from. All members must share the type and the declared size.

        - `type` **required** *string* &rarr;  RAID type: raid1 (case-insensitive).

        - `partition` optional &rarr;  Ordered partitions of the device; order determines placement on disk.
            - `fs` **required** *string* &rarr;  Filesystem type: ext4, xfs, vfat or swap (case-insensitive). Swap requires an empty mount.

            - `size` **required** *string* &rarr;  Partition size with an IEC suffix, for example 512MiB or 50GiB.

            - `mount` optional *string* &rarr;  Mount point of the partition. Empty or omitted means the partition is created but not mounted.

- `user_data` optional *string* &rarr;  Provide the cloud-init user-data payload.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the bare metal server.



## Import

A bare metal server can be imported using the `id`, e.g.
```shell
terraform import vkcs_baremetal_server.server 57b4d7f1-3acc-4843-b811-51df18badd9f
```
