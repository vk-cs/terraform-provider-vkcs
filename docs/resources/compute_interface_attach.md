---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_interface_attach"
description: |-
  Attaches a Network Interface to an Instance.
---

# vkcs_compute_interface_attach

Attaches a Network Interface (a Port) to an Instance using the VKCS Compute API.

## Example Usage
### Attachment Using an Existing Port
```terraform
resource "vkcs_compute_interface_attach" "etcd" {
  instance_id = vkcs_compute_instance.basic.id
  port_id     = vkcs_networking_port.persistent_etcd.id
}
```

### Attachment Using a Network ID
```terraform
resource "vkcs_compute_interface_attach" "db" {
  instance_id = vkcs_compute_instance.basic.id
  network_id  = vkcs_networking_network.db.id
}
```

Attaching multiple interfaces will not guarantee that interfaces are attached in
a deterministic manner. The interfaces will be attached in a seemingly random
order.
If you want to ensure that interfaces are attached in a given order, create
explicit dependencies between the interfaces , such as in virtual machines/vkcs_compute_volume_attach

## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.
```shell
terraform import vkcs_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
