---
subcategory: "VPN"
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_service"
description: |-
  Manages a VPN service resource within VKCS.
---

# vkcs_vpnaas_service

Manages a VPN service resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_service" "vpn_to_datacenter" {
  name = "vpn-tf-example"

  # See the argument description and check vkcs_networks_sdn datasource output to figure out
  # what type of router you should use in certain case (vkcs_networking_router or vkcs_dc_router)
  router_id = vkcs_dc_router.router.id

  depends_on = [vkcs_dc_interface.internet_sprut]
}
```
## Argument Reference
- `router_id` **required** *string* &rarr;  The ID of the router. Use router id for Neutron SDN and dc_router id for sprut SDN. To get a list of available SDNs in a project, you can use `vkcs_networking_sdn` datasource. Changing this creates a new service

- `admin_state_up` optional *boolean* &rarr;  The administrative state of the resource. Can either be up(true) or down(false). Changing this updates the administrative state of the existing service.

- `description` optional *string* &rarr;  The human-readable description for the service. Changing this updates the description of the existing service.

- `name` optional *string* &rarr;  The name of the service. Changing this updates the name of the existing service.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.<br>**New since v0.5.3**.

- `subnet_id` optional *string* &rarr;  SubnetID is the ID of the subnet. Default is null.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `external_v4_ip` *string* &rarr;  The read-only external (public) IPv4 address that is used for the VPN service.

- `external_v6_ip` *string* &rarr;  The read-only external (public) IPv6 address that is used for the VPN service.

- `id` *string* &rarr;  ID of the resource.

- `status` *string* &rarr;  Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_service.service_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
