---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_volume_attach"
description: |-
  Attaches a Block Storage Volume to an Instance.
---

# vkcs_compute_volume_attach

Attaches a Block Storage Volume to an Instance using the VKCS Compute API.

## Examples
### Usage with one volume
```terraform
resource "vkcs_compute_volume_attach" "data" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.data.id
}
```

### Usage with ORDERED multiple volumes
Attaching multiple volumes will not guarantee that the volumes are attached in
a deterministic manner. The volumes will be attached in a seemingly random
order.

If you want to ensure that the volumes are attached in a given order, create
explicit dependencies between the volumes, such as:
```terraform
resource "vkcs_compute_volume_attach" "attach_1" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volumes.0.id
}

resource "vkcs_compute_volume_attach" "attach_2" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volumes.1.id

  depends_on = [vkcs_compute_volume_attach.attach_1]
}
```
## Argument Reference
- `instance_id` **required** *string* &rarr;  The ID of the Instance to attach the Volume to.

- `volume_id` **required** *string* &rarr;  The ID of the Volume to attach to an Instance.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. A Compute client is needed to create a volume attachment. If omitted, the `region` argument of the provider is used. Changing this creates a new volume attachment.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Volume Attachments can be imported using the Instance ID and Volume ID separated by a slash, e.g.

```shell
terraform import vkcs_compute_volume_attach.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
