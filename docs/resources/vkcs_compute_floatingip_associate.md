---
layout: "vkcs"
page_title: "vkcs: vkcs_compute_floatingip_associate"
description: |-
  Associate a floating IP to an instance
---

# vkcs_compute_floatingip_associate

Associate a floating IP to an instance.

## Example Usage
### Automatically detect the correct network
```terraform
resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "my_pool"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
```

### Explicitly set the network to attach to
```terraform
resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }

  network {
    name = "default"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "my_pool"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  fixed_ip    = "${vkcs_compute_instance.instance_1.network.1.fixed_ip_v4}"
}
```
## Argument Reference
- `floating_ip` **String** (***Required***) The floating IP to associate.

- `instance_id` **String** (***Required***) The instance to associate the floating IP with.

- `fixed_ip` **String** (*Optional*) The specific IP address to direct traffic to.

- `region` **String** (*Optional*) The region in which to obtain the V2 Compute client. Keypairs are associated with accounts, but a Compute client is needed to create one. If omitted, the `region` argument of the provider is used. Changing this creates a new floatingip_associate.

- `wait_until_associated` **Boolean** (*Optional*) In cases where the VKCS environment does not automatically wait until the association has finished, set this option to have Terraform poll the instance until the floating IP has been associated. Defaults to false.


## Attributes Reference
- `floating_ip` **String** See Argument Reference above.

- `instance_id` **String** See Argument Reference above.

- `fixed_ip` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `wait_until_associated` **Boolean** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

This resource can be imported by specifying all three arguments, separated by a forward slash:
```shell
terraform import vkcs_compute_floatingip_associate.fip_1 <floating_ip>/<instance_id>/<fixed_ip>
```
