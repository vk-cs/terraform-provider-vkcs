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
  name = "basic-tf-example"
  # AZ and flavor are mandatory
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  # Use block_device to specify instance disk to get full control
  # of it in the future
  block_device {
    source_type      = "image"
    uuid             = data.vkcs_images_image.debian.id
    destination_type = "volume"
    volume_size      = 10
    volume_type      = "ceph-ssd"
    # Must be set to delete volume after instance deletion
    # Otherwise you get "orphaned" volume with terraform
    delete_on_termination = true
  }
  # Specify at least one network to not depend on project assets
  network {
    uuid = vkcs_networking_network.app.id
  }
  # Specify required security groups if you do not want `default` one
  security_group_ids = [
    vkcs_networking_secgroup.admin.id
  ]

  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
```
Use `vkcs_compute_floatingip_associate` to make the instance accessible from Internet.

### Instance with volume, tags and external IP
~> **Attention:** First, you should create the block storage volume and then attach it to the instance. Failing to do so will result in the virtual machine being provisioned with an ephemeral disk instead. Ephemeral disks lack certain capabilities, such as the ability to move or resize them. It's essential to adhere to the correct order of operations to avoid limitations in the management of block storage.
```terraform
resource "vkcs_compute_instance" "volumes_tags_externalip" {
  name              = "volumes-tags-externalip-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  # Use previously created volume as root device
  block_device {
    # Set boot_index to mark root device if multiple
    # block devices are specified
    boot_index       = 0
    source_type      = "volume"
    uuid             = vkcs_blockstorage_volume.bootable.id
    destination_type = "volume"
    # Omitting delete_on_termination (or setting it to false)
    # allows you to manage previously created volume after instance deletion
  }
  # Gracefully shutdown instance before deleting to keep
  # data consistent on persistent volume
  stop_before_destroy = true
  # Add empty disk to use ii during the instance lifecycle
  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 20
    delete_on_termination = true
  }
  tags = ["tf-example"]
  # Create the instance in external network to get external IP automatically
  # instead of using FIP on private network
  network {
    name = "ext-net"
  }
}
```

### Instance with Multiple Networks
~> **Attention:** When an instance connects to a network using an existing port, the instance's security groups are not automatically applied to the port. Instead, the security groups that were initially defined when the port was created will be applied to the port. As a result, the instance's security groups will not take effect in this scenario.
```terraform
resource "vkcs_compute_instance" "multiple_networks" {
  name              = "multiple-networks-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  # Autocreate a new port in 'app' network
  network {
    uuid = vkcs_networking_network.app.id
  }
  # Use previously created port
  # This does not change security groups associated with the port
  # Also this changes DNS name of the port
  network {
    port = vkcs_networking_port.persistent_etcd.id
  }
  # Attach 'admin' security group to autocreated port
  # This does not associate the group to 'persistent' port
  security_group_ids = [
    vkcs_networking_secgroup.admin.id
  ]
  depends_on = [
    vkcs_networking_router_interface.app,
    vkcs_networking_router_interface.db
  ]
}
```

### Instance with default and custom security groups
```terraform
resource "vkcs_compute_instance" "front_worker" {
  name      = "front-worker-tf-example"
  flavor_id = data.vkcs_compute_flavor.basic.id

  block_device {
    source_type      = "image"
    uuid             = data.vkcs_images_image.debian.id
    destination_type = "volume"
    volume_size      = 10
    # Must be set to delete volume after instance deletion
    # Otherwise you get "orphaned" volume with terraform
    delete_on_termination = true
  }

  security_group_ids = [
    data.vkcs_networking_secgroup.default_secgroup.id,
    vkcs_networking_secgroup.admin.id,
    vkcs_networking_secgroup.http.id
  ]
  image_id = data.vkcs_images_image.debian.id

  network {
    uuid = vkcs_networking_network.app.id
  }
}
```

### Instance with personality
~> **Attention:** To use this feature, you must set the `config_drive` argument to true.
The Personality feature allows you to customize the files and scripts that are injected into an instance during its provisioning. When using the Personality feature, you can provide one or more files or scripts that are associated with an instance. During instance creation, these files are injected into the instance's file system.

This feature is useful for tasks such as bootstrapping instances with custom configurations, deploying specific software or packages, and executing initialization scripts.
```terraform
resource "vkcs_compute_instance" "basic" {
  name              = "personality-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.app.id
  }
  # config_drive must be enabled to use personality
  config_drive = true
  personality {
    file    = "/opt/app/config.json"
    content = jsonencode({ "foo" : "bar" })
  }
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
```

### Instance with user data and cloud-init
This feature is used to provide initialization scripts or configurations to instances during their launch. User data is typically used to automate the configuration and customization of instances at the time of provisioning.
```terraform
resource "vkcs_compute_instance" "basic" {
  name              = "basic-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.app.id
  }
  user_data = <<EOF
    #cloud-config
    package_upgrade: true
    packages:
      - nginx
    runcmd:
      - systemctl start nginx
  EOF
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
```
Also, the user_data option can be set to the contents of a script file using the file() function:
  user_data = file("${path.module}/userdata.sh")

### Instance with cloud monitoring

~> **Attention:** If you enable cloud monitoring and use user_data, the terraform provider will try to merge monitoring
script and user_data into one file

```terraform
resource "vkcs_compute_instance" "cloud_monitoring" {
  name              = "cloud-monitoring-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.app.id
  }

  cloud_monitoring {
    service_user_id = vkcs_cloud_monitoring.basic.service_user_id
    script          = vkcs_cloud_monitoring.basic.script
  }

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
```

## Argument Reference
- `name` **required** *string* &rarr;  A unique name for the resource.

- `access_ip_v4` optional *string* &rarr;  The first detected Fixed IPv4 address.

- `admin_pass` optional sensitive *string* &rarr;  The administrative password to assign to the server. Changing this changes the root password on the existing server.

- `availability_zone` optional *string* &rarr;  The availability zone in which to create the server. Conflicts with `availability_zone_hints`. Changing this creates a new server.

- `block_device` optional &rarr;  Configuration of block devices. The block_device structure is documented below. Changing this creates a new server. You can specify multiple block devices which will create an instance with multiple disks. This configuration is very flexible, so please see the following [reference](https://docs.openstack.org/nova/latest/user/block-device-mapping.html) for more information.
  - `source_type` **required** *string* &rarr;  The source type of the device. Must be one of "blank", "image", "volume", or "snapshot". Changing this creates a new server.

  - `boot_index` optional *number* &rarr;  The boot index of the volume. It defaults to 0 if only one `block_device` is specified, and to -1 if more than one is configured. Changing this creates a new server. <br>**Note:** You must set the boot index to 0 for one of the block devices if more than one is defined.

  - `delete_on_termination` optional *boolean* &rarr;  Delete the volume / block device upon termination of the instance. Defaults to false. Changing this creates a new server. _<br>**Note:**_ It is important to enable `delete_on_termination` for volumes created with instance. If `delete_on_termination` is disabled for such volumes, then after instance deletion such volumes will stay orphaned and uncontrolled by terraform. _<br>**Note:**_ It is important to disable `delete_on_termination` if volume is created as separate terraform resource and is attached to instance. Enabling `delete_on_termination` for such volumes will result in mismanagement between two terraform resources in case of instance deletion

  - `destination_type` optional *string* &rarr;  The type that gets created. Possible values are "volume" and "local". Changing this creates a new server.

  - `device_type` optional *string* &rarr;  The low-level device type that will be used. Most common thing is to leave this empty. Changing this creates a new server.

  - `disk_bus` optional *string* &rarr;  The low-level disk bus that will be used. Most common thing is to leave this empty. Changing this creates a new server.

  - `guest_format` optional *string* &rarr;  Specifies the guest server disk file system format, such as `ext2`, `ext3`, `ext4`, `xfs` or `swap`. Swap block device mappings have the following restrictions: source_type must be blank and destination_type must be local and only one swap disk per server and the size of the swap disk must be less than or equal to the swap size of the flavor. Changing this creates a new server.

  - `uuid` optional *string* &rarr;  The UUID of the image, volume, or snapshot. Optional if `source_type` is set to `"blank"`. Changing this creates a new server.

  - `volume_size` optional *number* &rarr;  The size of the volume to create (in gigabytes). Required in the following combinations: source=image and destination=volume, source=blank and destination=local, and source=blank and destination=volume. Changing this creates a new server.

  - `volume_type` optional *string* &rarr;  The volume type that will be used. Changing this creates a new server.

- `cloud_monitoring` optional &rarr;  The settings of the cloud monitoring, it is recommended to set this field with the values of `vkcs_cloud_monitoring` resource fields. Changing this creates a new server.
  - `script` **required** sensitive *string* &rarr;  The script of the cloud monitoring.

  - `service_user_id` **required** *string* &rarr;  The id of the service monitoring user.

- `config_drive` optional *boolean* &rarr;  Whether to use the config_drive feature to configure the instance. Changing this creates a new server.

- `flavor_id` optional *string* &rarr;  The flavor ID of the desired flavor for the server. Required if `flavor_name` is empty. Changing this resizes the existing server.

- `flavor_name` optional *string* &rarr;  The name of the desired flavor for the server. Required if `flavor_id` is empty. Changing this resizes the existing server.

- `force_delete` optional *boolean* &rarr;  Whether to force the compute instance to be forcefully deleted. This is useful for environments that have reclaim / soft deletion enabled.

- `image_id` optional *string* &rarr;  The image ID of the desired image for the server. Required if `image_name` is empty and not booting from a volume. Do not specify if booting from a volume. Changing this creates a new server.

- `image_name` optional *string* &rarr;  The name of the desired image for the server. Required if `image_id` is empty and not booting from a volume. Do not specify if booting from a volume. Changing this creates a new server.

- `key_pair` optional *string* &rarr;  The name of a key pair to put on the server. The key pair must already be created and associated with the tenant's account. Changing this creates a new server.

- `metadata` optional *map of* *string* &rarr;  Metadata key/value pairs to make available from within the instance. Changing this updates the existing server metadata.

- `network` optional &rarr;  An array of one or more networks to attach to the instance. The network object structure is documented below. Changing this creates a new server.
  - `access_network` optional *boolean* &rarr;  Specifies if this network should be used for provisioning access. Accepts true or false. Defaults to false.

  - `fixed_ip_v4` optional *string* &rarr;  Specifies a fixed IPv4 address to be used on this network. Changing this creates a new server.

  - `name` optional *string* &rarr;  The human-readable name of the network. Optional if `uuid` or `port` is provided. Changing this creates a new server.

  - `port` optional *string* &rarr;  The port UUID of a network to attach to the server. Optional if `uuid` or `name` is provided. Changing this creates a new server. <br>**Note:** If port is used, only its security groups will be applied instead of security_groups instance argument.

  - `uuid` optional *string* &rarr;  The network UUID to attach to the server. Optional if `port` or `name` is provided. Changing this creates a new server.

- `network_mode` optional *string* &rarr;  Special string for `network` option to create the server. `network_mode` can be `"auto"` or `"none"`. Please see the following [reference](https://docs.openstack.org/api-ref/compute/?expanded=create-server-detail#id11) for more information. Conflicts with `network`.

- `personality` optional &rarr;  Customize the personality of an instance by defining one or more files and their contents. The personality structure is described below. <br>**Note:** 'config_drive' must be enabled.
  - `content` **required** *string* &rarr;  The contents of the file.

  - `file` **required** *string* &rarr;  The absolute path of the destination file. Limited to 255 bytes.

- `power_state` optional *string* &rarr;  Provide the VM state. Only 'active' and 'shutoff' are supported values. <br>**Note:** If the initial power_state is the shutoff the VM will be stopped immediately after build and the provisioners like remote-exec or files are not supported.

- `region` optional *string* &rarr;  The region in which to create the server instance. If omitted, the `region` argument of the provider is used. Changing this creates a new server.

- `scheduler_hints` optional &rarr;  Provide the Nova scheduler with hints on how the instance should be launched. The available hints are described below.
  - `group` optional *string* &rarr;  A UUID of a Server Group. The instance will be placed into that group.

- `security_group_ids` optional *set of* *string* &rarr;  An array of one or more security group ids to associate with the server. Changing this results in adding/removing security groups from the existing server. <br>**Note:** When attaching the instance to networks using Ports, place the security groups on the Port and not the instance.<br>**New since v0.7.3**.

- `security_groups` optional deprecated *set of* *string* &rarr;  An array of one or more security group names to associate with the server. Changing this results in adding/removing security groups from the existing server. <br>**Note:** When attaching the instance to networks using Ports, place the security groups on the Port and not the instance. **Deprecated** Configure `security_group_ids` instead.

- `stop_before_destroy` optional *boolean* &rarr;  Whether to try stop instance gracefully before destroying it, thus giving chance for guest OS daemons to stop correctly. If instance doesn't stop within timeout, it will be destroyed anyway.

- `tags` optional *set of* *string* &rarr;  A set of string tags for the instance. Changing this updates the existing instance tags.

- `user_data` optional *string* &rarr;  The user data to provide when launching the instance. When cloud_monitoring enabled only #!/bin/bash, #cloud-config, #ps1 user_data formats are supported. Changing this creates a new server.

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
### Instances and network
* When creating a network for the instance and connecting it to the router, please ensure that you specify the dependency of the instance on the `vkcs_networking_router_interface` resource. This is crucial for the proper initialization of the instance in such scenarios.
* An instance cannot be created without a valid network configuration even if you intend to use `vkcs_compute_interface_attach` after the instance has been created. Please note that if the network block is not explicitly specified, it will be automatically created implicitly only if there is a single network configuration in the project. However, if there are multiple network configurations, the process will fail.

### Instances and Security Groups
~> **Attention:** If you specify the ID of the security group instead of name in old argument _security_groups_, terraform will remove and reapply the security group upon each call.

When referencing a security group resource in an instance resource, always use _security_group_id_ argument.
If you want replace old argument _security_groups_ with _security_group_ids_ just find their ID's and change resource in terraform file, your cloud resource will not be changed.

```hcl
resource "vkcs_networking_secgroup" "sg_1" {
  name = "sg_1"
}

resource "vkcs_compute_instance" "foo" {
  name            = "terraform-test"
  security_group_ids = [vkcs_networking_secgroup.sg_1.id]
}
```

### Tags and metadata
**Metadata:** The `metadata` option allows you to attach key-value pairs of metadata to an instance. Metadata provides additional contextual information about the instance. You could include metadata specifying custom settings, application-specific data, or any other parameters needed for the instance's configuration or operation. Also, metadata can be leveraged by external systems or tools to automate actions or perform integrations. For example, external systems can read instance metadata to retrieve specific details or trigger actions based on the metadata values.
```hcl
resource "vkcs_compute_instance" "instance_1" {

  # Other instance configuration...

  metadata = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

**Tags:** The `tags` option allows you to assign one or more labels (tags) to an instance. Tags are a way to categorize and organize instances based on specific attributes. Tags can be used as triggers for automation or policy enforcement. For example, you might have automation scripts or policies that automatically perform specific actions on instances based on their tags.
```hcl
resource "vkcs_compute_instance" "instance_1" {

  # Other instance configuration...

  tags = ["webserver", "production"]
}
```

By using the `metadata` and `tags` options, you can provide additional context, organization, and categorization to instances, making it easier to manage and identify them based on specific attributes or requirements.

### User-data and personality
**User-data:** The `user-data` option is used to provide cloud-init configuration data to an instance, while the "personality" option is used to inject files into an instance during its creation. Cloud-init is a widely used multi-distribution package that handles early initialization of a cloud instance.

**Personality:** The `personality` option allows you to inject files into an instance during its creation. You can provide one or more files that will be available inside the instance's file system. This can be useful for provisioning additional configuration files, scripts, or any other required data.

Both options, `user-data` and `personality` can be used in combination or individually to customize instances. The choice between them depends on your specific needs and whether you want to provide configuration data through cloud-init or inject files directly into the instance's file system.

## Importing instances

Importing instances can be tricky, since the nova api does not offer all information provided at creation time for later retrieval.
Compute returns the network interfaces grouped by network, thus not in creation order. That means that if you have multiple network interfaces you must take care of the order of networks in your configuration, or read the network order in the state file after import and modify your configuration accordingly.

~> **Note:**
A note on ports. If you have created a networking port independent of an instance, then the import code has no way to detect that the port is created idenpendently, and therefore on deletion of imported instances you might have port resources in your project, which you expected to be created by the instance and thus to also be deleted with the instance.

~> **Note:**
A note on block storage volumes, the importer does not read `delete_on_termination` flag, and always assumes true. If you import an instance created with `delete_on_termination` false, you end up with "orphaned" volumes after destruction of instances.

To import instance use following command:
```shell
terraform import vkcs_compute_instance.basic_instance b61e8c9a-94ca-4852-9008-a95cdae6a2d9
```
