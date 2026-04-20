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

### Attachment Using a Network ID with a Fixed IP
```terraform
resource "vkcs_compute_interface_attach" "db" {
  instance_id = vkcs_compute_instance.basic.id
  network_id  = vkcs_networking_network.db.id
  fixed_ip    = "192.0.2.11"
}
```

Attaching multiple interfaces will not guarantee that interfaces are attached in
a deterministic manner. The interfaces will be attached in a seemingly random
order.
If you want to ensure that interfaces are attached in a given order, create
explicit dependencies between the interfaces , such as in virtual machines/vkcs_compute_volume_attach

## Argument Reference
- `instance_id` **required** *string* &rarr; The ID of the Instance to attach the Port or Network to.

- `region` _optional_ *string* &rarr;  The region in which to create the interface attachment. <br>If omitted, the `region` argument of the provider is used. Changing this creates a new attachment.

- `port_id` _optional_ *string* &rarr; The ID of the Port to attach to an Instance. <br>_Note_: This option and `network_id` are mutually exclusive.

- `network_id` _optional_ *string* &rarr; The ID of the Network to attach to an Instance. A port will be created automatically. <br>_Note_: This option and `port_id` are mutually exclusive.

- `fixed_ip` _optional_ *string* &rarr; An IP address to assosciate with the port. <br>You must specify a `network_id`, the IP address must lie in a range on the supplied network. <br>_Note_: This option cannot be used with `port_id`.

## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.
```shell
terraform import vkcs_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
