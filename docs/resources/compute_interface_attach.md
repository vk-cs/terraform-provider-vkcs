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
### Basic Attachment
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  network_id  = "${vkcs_networking_port.network_1.id}"
}
```

### Attachment Specifying a Fixed IP
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  network_id  = "${vkcs_networking_port.network_1.id}"
  fixed_ip    = "10.0.10.10"
}
```

### Attachment Using an Existing Port
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_port" "port_1" {
  name           = "port_1"
  network_id     = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"
}


resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  port_id     = "${vkcs_networking_port.port_1.id}"
}
```

### Attaching Multiple Interfaces
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_port" "ports" {
  count          = 2
  name           = "${format("port-%02d", count.index + 1)}"
  network_id     = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "attachments" {
  count       = 2
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  port_id     = "${vkcs_networking_port.ports.*.id[count.index]}"
}
```

Note that the above example will not guarantee that the ports are attached in a deterministic manner. The ports will be attached in a seemingly random order.

If you want to ensure that the ports are attached in a given order, create explicit dependencies between the ports, such as:
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_port" "ports" {
  count          = 2
  name           = "${format("port-%02d", count.index + 1)}"
  network_id     = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"
}

resource "vkcs_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  port_id     = "${vkcs_networking_port.ports.*.id[0]}"
}

resource "vkcs_compute_interface_attach" "ai_2" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  port_id     = "${vkcs_networking_port.ports.*.id[1]}"
}
```
## Argument Reference
- `instance_id` **required** *string* &rarr;  The ID of the Instance to attach the Port or Network to.

- `fixed_ip` optional *string* &rarr;  An IP address to assosciate with the port.
_NOTE_: This option cannot be used with port_id. You must specify a network_id. The IP address must lie in a range on the supplied network.

- `network_id` optional *string* &rarr;  The ID of the Network to attach to an Instance. A port will be created automatically.
_NOTE_: This option and `port_id` are mutually exclusive.

- `port_id` optional *string* &rarr;  The ID of the Port to attach to an Instance.
_NOTE_: This option and `network_id` are mutually exclusive.

- `region` optional *string* &rarr;  The region in which to create the interface attachment. If omitted, the `region` argument of the provider is used. Changing this creates a new attachment.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.
```shell
terraform import vkcs_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
