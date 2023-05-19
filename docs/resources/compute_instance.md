---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_instance"
description: |-
  Manages a compute VM instance.
---

# vkcs_compute_instance

Manages a compute VM instance resource.

~> **Note:** All arguments including the instance admin password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage
### Basic Instance
```terraform
resource "vkcs_compute_instance" "basic" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  metadata = {
    this = "that"
  }

  network {
    name = "my_network"
  }
}
```

### Instance With Attached Volume
```terraform
resource "vkcs_blockstorage_volume" "myvol" {
  name = "myvol"
  size = 1
}

resource "vkcs_compute_instance" "myinstance" {
  name            = "myinstance"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }
}

resource "vkcs_compute_volume_attach" "attached" {
  instance_id = "${vkcs_compute_instance.myinstance.id}"
  volume_id   = "${vkcs_blockstorage_volume.myvol.id}"
}
```

### Boot From Volume
```terraform
resource "vkcs_compute_instance" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    volume_size           = 5
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
```

### Boot From an Existing Volume
```terraform
resource "vkcs_blockstorage_volume" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = "<image-id>"
}

resource "vkcs_compute_instance" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "${vkcs_blockstorage_volume.myvol.id}"
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
```

### Boot Instance, Create Volume, and Attach Volume as a Block Device
```terraform
resource "vkcs_compute_instance" "compute" {
  name            = "compute-instance"
  flavor_id       = data.vkcs_compute_flavor.compute.id
  security_groups = ["default"]
  availability_zone = "GZ1"

  block_device {
    uuid                  = data.vkcs_images_image.compute.id
    source_type           = "image"
    destination_type      = "volume"
    volume_type           = "ceph-ssd"
    volume_size           = 8
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_type           = "ceph-ssd"
    volume_size           = 8
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.compute.id
  }

  depends_on = [
    vkcs_networking_network.compute,
    vkcs_networking_subnet.compute
  ]
}

resource "vkcs_networking_floatingip" "fip" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip" {
  floating_ip = vkcs_networking_floatingip.fip.address
  instance_id = vkcs_compute_instance.compute.id
}
```

### Boot Instance and Attach Existing Volume as a Block Device
```terraform
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    uuid                  = "${vkcs_blockstorage_volume.volume_1.id}"
    source_type           = "volume"
    destination_type      = "volume"
    boot_index            = 1
    delete_on_termination = true
  }
  network {
    name = "my_network"
  }
}
```

### Instance With Multiple Networks
```terraform
resource "vkcs_networking_floatingip" "myip" {
  pool = "my_pool"
}

resource "vkcs_compute_instance" "multi-net" {
  name            = "multi-net"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_first_network"
  }

  network {
    name = "my_second_network"
  }
}

resource "vkcs_compute_floatingip_associate" "myip" {
  floating_ip = "${vkcs_networking_floatingip.myip.address}"
  instance_id = "${vkcs_compute_instance.multi-net.id}"
  fixed_ip    = "${vkcs_compute_instance.multi-net.network.1.fixed_ip_v4}"
}
```

### Instance With Personality
```terraform
resource "vkcs_compute_instance" "personality" {
  name            = "personality"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  personality {
    file    = "/path/to/file/on/instance.txt"
    content = "contents of file"
  }

  network {
    name = "my_network"
  }
}
```

### Instance with Multiple Ephemeral Disks
```terraform
resource "vkcs_compute_instance" "multi-eph" {
  name            = "multi_eph"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }
  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
    guest_format          = "ext4"
  }
  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
  network {
    name = "my_network"
  }
}
```

### Instance with Boot Disk and Swap Disk
```terraform
resource "vkcs_compute_flavor" "flavor-with-swap" {
  name  = "flavor-with-swap"
  ram   = "8096"
  vcpus = "2"
  disk  = "20"
  swap  = "4096"
}
resource "vkcs_compute_instance" "vm-swap" {
  name            = "vm_swap"
  flavor_id       = "${vkcs_compute_flavor.flavor-with-swap.id}"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }
  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    guest_format          = "swap"
    volume_size           = 4
  }
  network {
    name = "my_network"
  }
}
```

### Instance with User Data (cloud-init)
```terraform
resource "vkcs_compute_instance" "instance_1" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  user_data       = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    name = "my_network"
  }
}
```
`user_data` can come from a variety of sources: inline, read in from the `file` function, or the `template_cloudinit_config` resource.
## Argument Reference
- `name` **required** *string* &rarr;  A unique name for the resource.

- `access_ip_v4` optional *string* &rarr;  The first detected Fixed IPv4 address.

- `admin_pass` optional sensitive *string* &rarr;  The administrative password to assign to the server. Changing this changes the root password on the existing server.

- `availability_zone` optional *string* &rarr;  The availability zone in which to create the server. Conflicts with `availability_zone_hints`. Changing this creates a new server.

- `block_device` optional &rarr;  Configuration of block devices. The block_device structure is documented below. Changing this creates a new server. You can specify multiple block devices which will create an instance with multiple disks. This configuration is very flexible, so please see the following [reference](https://docs.openstack.org/nova/latest/user/block-device-mapping.html) for more information.
  - `source_type` **required** *string* &rarr;  The source type of the device. Must be one of "blank", "image", "volume", or "snapshot". Changing this creates a new server.

  - `boot_index` optional *number* &rarr;  The boot index of the volume. It defaults to -1. Changing this creates a new server. <br>**Note:** You must set the boot index to 0 for one of the block devices if more than one is defined.

  - `delete_on_termination` optional *boolean* &rarr;  Delete the volume / block device upon termination of the instance. Defaults to false. Changing this creates a new server.

  - `destination_type` optional *string* &rarr;  The type that gets created. Possible values are "volume" and "local". Changing this creates a new server.

  - `device_type` optional *string* &rarr;  The low-level device type that will be used. Most common thing is to leave this empty. Changing this creates a new server.

  - `disk_bus` optional *string* &rarr;  The low-level disk bus that will be used. Most common thing is to leave this empty. Changing this creates a new server.

  - `guest_format` optional *string* &rarr;  Specifies the guest server disk file system format, such as `ext2`, `ext3`, `ext4`, `xfs` or `swap`. Swap block device mappings have the following restrictions: source_type must be blank and destination_type must be local and only one swap disk per server and the size of the swap disk must be less than or equal to the swap size of the flavor. Changing this creates a new server.

  - `uuid` optional *string* &rarr;  (Optional if `source_type` is set to `"blank"` ) The UUID of the image, volume, or snapshot. Changing this creates a new server.

  - `volume_size` optional *number* &rarr;  The size of the volume to create (in gigabytes). Required in the following combinations: source=image and destination=volume, source=blank and destination=local, and source=blank and destination=volume. Changing this creates a new server.

  - `volume_type` optional *string* &rarr;  The volume type that will be used. Changing this creates a new server.

- `config_drive` optional *boolean* &rarr;  Whether to use the config_drive feature to configure the instance. Changing this creates a new server.

- `flavor_id` optional *string* &rarr;  (Required if `flavor_name` is empty) The flavor ID of the desired flavor for the server. Changing this resizes the existing server.

- `flavor_name` optional *string* &rarr;  (Required if `flavor_id` is empty) The name of the desired flavor for the server. Changing this resizes the existing server.

- `force_delete` optional *boolean* &rarr;  Whether to force the compute instance to be forcefully deleted. This is useful for environments that have reclaim / soft deletion enabled.

- `image_id` optional *string* &rarr;  (Required if `image_name` is empty and not booting from a volume. Do not specify if booting from a volume.) The image ID of the desired image for the server. Changing this creates a new server.

- `image_name` optional *string* &rarr;  (Required if `image_id` is empty and not booting from a volume. Do not specify if booting from a volume.) The name of the desired image for the server. Changing this creates a new server.

- `key_pair` optional *string* &rarr;  The name of a key pair to put on the server. The key pair must already be created and associated with the tenant's account. Changing this creates a new server.

- `metadata` optional *map of* *string* &rarr;  Metadata key/value pairs to make available from within the instance. Changing this updates the existing server metadata.

- `network` optional &rarr;  An array of one or more networks to attach to the instance. The network object structure is documented below. Changing this creates a new server.
  - `access_network` optional *boolean* &rarr;  Specifies if this network should be used for provisioning access. Accepts true or false. Defaults to false.

  - `fixed_ip_v4` optional *string* &rarr;  Specifies a fixed IPv4 address to be used on this network. Changing this creates a new server.

  - `name` optional *string* &rarr;  (Optional if `uuid` or `port` is provided) The human-readable name of the network. Changing this creates a new server.

  - `port` optional *string* &rarr;  (Optional if `uuid` or `name` is provided) The port UUID of a network to attach to the server. Changing this creates a new server.

  - `uuid` optional *string* &rarr;  (Optional if `port`  or `name` is provided) The network UUID to attach to the server. Changing this creates a new server.

- `network_mode` optional *string* &rarr;  Special string for `network` option to create the server. `network_mode` can be `"auto"` or `"none"`. Please see the following [reference](https://docs.openstack.org/api-ref/compute/?expanded=create-server-detail#id11) for more information. Conflicts with `network`.

- `personality` optional &rarr;  Customize the personality of an instance by defining one or more files and their contents. The personality structure is described below.
  - `content` **required** *string* &rarr;  The contents of the file. Limited to 255 bytes.

  - `file` **required** *string* &rarr;  The absolute path of the destination file.

- `power_state` optional *string* &rarr;  Provide the VM state. Only 'active' and 'shutoff' are supported values. *Note*: If the initial power_state is the shutoff the VM will be stopped immediately after build and the provisioners like remote-exec or files are not supported.

- `region` optional *string* &rarr;  The region in which to create the server instance. If omitted, the `region` argument of the provider is used. Changing this creates a new server.

- `scheduler_hints` optional &rarr;  Provide the Nova scheduler with hints on how the instance should be launched. The available hints are described below.
  - `group` optional *string* &rarr;  A UUID of a Server Group. The instance will be placed into that group.

- `security_groups` optional *set of* *string* &rarr;  An array of one or more security group names to associate with the server. Changing this results in adding/removing security groups from the existing server. *Note*: When attaching the instance to networks using Ports, place the security groups on the Port and not the instance. *Note*: Names should be used and not ids, as ids trigger unnecessary updates.

- `stop_before_destroy` optional *boolean* &rarr;  Whether to try stop instance gracefully before destroying it, thus giving chance for guest OS daemons to stop correctly. If instance doesn't stop within timeout, it will be destroyed anyway.

- `tags` optional *set of* *string* &rarr;  A set of string tags for the instance. Changing this updates the existing instance tags.

- `user_data` optional *string* &rarr;  The user data to provide when launching the instance.	Changing this creates a new server.

- `vendor_options` optional &rarr;  Map of additional vendor-specific options. Supported options are described below.
  - `detach_ports_before_destroy` optional *boolean* &rarr;  Whether to try to detach all attached ports to the vm before destroying it to make sure the port state is correct after the vm destruction. This is helpful when the port is not deleted.

  - `ignore_resize_confirmation` optional *boolean* &rarr;  Boolean to control whether to ignore manual confirmation of the instance resizing.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_metadata` *map of* *string* &rarr;  Contains all instance metadata, even metadata not set by Terraform.

- `all_tags` *set of* *string* &rarr;  The collection of tags assigned on the instance, which have been explicitly and implicitly added.

- `id` *string* &rarr;  ID of the resource.

- `network` 
  - `mac` *string* &rarr;  The MAC address of the NIC on that network.



## Notes
### Instances and Security Groups

When referencing a security group resource in an instance resource, always use the _name_ of the security group. If you specify the ID of the security group, Terraform will remove and reapply the security group upon each call. This is because the VKCS Compute API returns the names of the associated security groups and not their IDs.

Note the following example:

```hcl
resource "vkcs_networking_secgroup" "sg_1" {
  name = "sg_1"
}

resource "vkcs_compute_instance" "foo" {
  name            = "terraform-test"
  security_groups = ["${vkcs_networking_secgroup.sg_1.name}"]
}
```

### Instances and Ports

Neutron Ports are a great feature and provide a lot of functionality. However, there are some notes to be aware of when mixing Instances and Ports:

* When attaching an Instance to one or more networks using Ports, place the security groups on the Port and not the Instance. If you place the security groups on the Instance, the security groups will not be applied upon creation, but they will be applied upon a refresh.

* Network IP information is not available within an instance for networks that are attached with Ports. This is mostly due to the flexibility Neutron Ports provide when it comes to IP addresses. For example, a Neutron Port can have multiple Fixed IP addresses associated with it. It's not possible to know which single IP address the user would want returned to the Instance's state information. Therefore, in order for a Provisioner to connect to an Instance via it's network Port, customize the `connection` information:

```hcl
resource "vkcs_networking_port" "port_1" {
  name           = "port_1"
  admin_state_up = "true"

  network_id = "0a1d0a27-cffa-4de3-92c5-9d3fd3f2e74d"

  security_group_ids = [
    "2f02d20a-8dca-49b7-b26f-b6ce9fddaf4f",
    "ca1e5ed7-dae8-4605-987b-fadaeeb30461",
  ]
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"

  network {
    port = "${vkcs_networking_port.port_1.id}"
  }

  connection {
    user        = "root"
    host        = "${vkcs_networking_port.port_1.fixed_ip.0.ip_address}"
    private_key = "~/path/to/key"
  }

  provisioner "remote-exec" {
    inline = [
      "echo terraform executed > /tmp/foo",
    ]
  }
}
```

### Instances and Networks

Instances almost always require a network. Here are some notes to be aware of with how Instances and Networks relate:

* In scenarios where you only have one network available, you can create an instance without specifying a `network` block. VKCS will automatically launch the instance on this network.

* If you have access to more than one network, you will need to specify a network with a `network` block. Not specifying a network will result in the following error:

```
* vkcs_compute_instance.instance: Error creating VKCS server:
Expected HTTP response code [201 202] when accessing [POST https://example.com:8774/v2.1/servers], but got 409 instead
{"conflictingRequest": {"message": "Multiple possible networks found, use a Network ID to be more specific.", "code": 409}}
```

* If you intend to use the `vkcs_compute_interface_attach` resource, you still need to make sure one of the above points is satisfied. An instance cannot be created without a valid network configuration even if you intend to use `vkcs_compute_interface_attach` after the instance has been created.

## Importing instances

Importing instances can be tricky, since the nova api does not offer all information provided at creation time for later retrieval. Network interface attachment order, and number and sizes of ephemeral disks are examples of this.

### Importing basic instance
Assume you want to import an instance with one ephemeral root disk, and one network interface.

Your configuration would look like the following:

```hcl
resource "vkcs_compute_instance" "basic_instance" {
  name            = "basic"
  flavor_id       = "<flavor_id>"
  key_pair        = "<keyname>"
  security_groups = ["default"]
  image_id =  "<image_id>"

  network {
    name = "<network_name>"
  }
}

```
Then you execute
```shell
terraform import vkcs_compute_instance.basic_instance b61e8c9a-94ca-4852-9008-a95cdae6a2d9
```

### Importing instance with multiple network interfaces.

Compute returns the network interfaces grouped by network, thus not in creation order.
That means that if you have multiple network interfaces you must take care of the order of networks in your configuration.


As example we want to import an instance with one ephemeral root disk, and 3 network interfaces.

Examples

```hcl
resource "vkcs_compute_instance" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "<flavor_id"
  key_pair        = "<keyname>"
  image_id        = <image_id>
  security_groups = ["default"]

  network {
    name = "<network1>"
  }
  network {
    name = "<network2>"
  }
  network {
    name = "<network1>"
    fixed_ip_v4 = "<fixed_ip_v4>"
  }

}
```

In the above configuration the networks are out of order compared to what nova and thus the import code returns, which means the plan will not be empty after import.

So either with care check the plan and modify configuration, or read the network order in the state file after import and modify your configuration accordingly.

 * A note on ports. If you have created a networking port independent of an instance, then the import code has no way to detect that the port is created idenpendently, and therefore on deletion of imported instances you might have port resources in your project, which you expected to be created by the instance and thus to also be deleted with the instance.

### Importing instances with multiple block storage volumes.

We have an instance with two block storage volumes, one bootable and one non-bootable.
Note that we only configure the bootable device as block_device.
The other volumes can be specified as `vkcs_blockstorage_volume`

```hcl
resource "vkcs_compute_instance" "instance_2" {
  name            = "instance_2"
  image_id        = "<image_id>"
  flavor_id       = "<flavor_id>"
  key_pair        = "<keyname>"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image_id>"
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

   network {
    name = "<network_name>"
  }
}
resource "vkcs_blockstorage_volume" "volume_1" {
  size = 1
  name = "<vol_name>"
}
resource "vkcs_compute_volume_attach" "va_1" {
  volume_id   = "${vkcs_blockstorage_volume.volume_1.id}"
  instance_id = "${vkcs_compute_instance.instance_2.id}"
}
```
To import the instance outlined in the above configuration do the following:

```shell
terraform import vkcs_compute_instance.instance_2 9758f738-3ef3-4c58-aedf-9e0378855ff2
terraform import vkcs_blockstorage_volume.volume_1 47ae2945-0d92-47a2-8bbb-bc007f161272
terraform import vkcs_compute_volume_attach.va_1 9758f738-3ef3-4c58-aedf-9e0378855ff2/47ae2945-0d92-47a2-8bbb-bc007f161272
```

* A note on block storage volumes, the importer does not read delete_on_termination flag, and always assumes true. If you import an instance created with delete_on_termination false, you end up with "orphaned" volumes after destruction of instances.
