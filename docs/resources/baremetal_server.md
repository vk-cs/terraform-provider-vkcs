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
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

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
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

  bond {
    name = "bond0"
    interface_names = ["nic0", "nic1"]
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }

    vlan {
      id = 100
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
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

  bond {
    name = "bond0"
    interface_names = ["nic0"]
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }
  }

  bond {
    name = "bond1"
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
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

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

- `nic` optional &rarr;  Physical network interfaces.
    - `name` **required** *string* &rarr;  Interface name (e.g. nic0, eno1). Acts as unique identifier.

    - `vlan` optional &rarr;  VLAN configuration. Allowed only if interface is not part of a bond.
        - `id` optional *number* &rarr;  Number of the VLAN.

        - `native` optional *boolean* &rarr;  Whether the VLAN is native.

        - `network_id` optional *string* &rarr;  ID of the network.

        - `subnet_id` optional *string* &rarr;  ID of the subnet.

- `os_id` optional *string* &rarr;  Set os id.

- `raid_type` optional *string* &rarr;  Parameter to determine should RAID be used during image flashing.

- `region` optional *string* &rarr;  The region to fetch the bare metal server from, defaults to the provider's region.

- `user_data` optional *string* &rarr;  Provide the cloud-init user-data payload.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the bare metal server.



## Import

A bare metal server can be imported using the `id`, e.g.
```shell
terraform import vkcs_baremetal_server.server 57b4d7f1-3acc-4843-b811-51df18badd9f
```
