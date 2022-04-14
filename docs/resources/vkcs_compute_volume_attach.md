---
layout: "vkcs"
page_title: "vkcs: compute_volume_attach"
description: |-
  Attaches a Block Storage Volume to an Instance.
---

# vkcs\_compute\_volume\_attach

Attaches a Block Storage Volume to an Instance using the VKCS
Compute API.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volume_1.id}"
}
```

### Attaching multiple volumes to a single instance

```hcl
resource "vkcs_blockstorage_volume" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "attachments" {
  count       = 2
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.*.id[count.index]}"
}

output "volume_devices" {
  value = "${vkcs_compute_volume_attach.attachments.*.device}"
}
```

Note that the above example will not guarantee that the volumes are attached in
a deterministic manner. The volumes will be attached in a seemingly random
order.

If you want to ensure that the volumes are attached in a given order, create
explicit dependencies between the volumes, such as:

```hcl
resource "vkcs_blockstorage_volume" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "attach_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.0.id}"
}

resource "vkcs_compute_volume_attach" "attach_2" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.1.id}"

  depends_on = ["vkcs_compute_volume_attach.attach_1"]
}

output "volume_devices" {
  value = "${vkcs_compute_volume_attach.attachments.*.device}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Compute client.
    A Compute client is needed to create a volume attachment. If omitted, the
    `region` argument of the provider is used. Changing this creates a
    new volume attachment.

* `instance_id` - (Required) The ID of the Instance to attach the Volume to.

* `volume_id` - (Required) The ID of the Volume to attach to an Instance.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `volume_id` - See Argument Reference above.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID
separated by a slash, e.g.

```
$ terraform import vkcs_compute_volume_attach.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
