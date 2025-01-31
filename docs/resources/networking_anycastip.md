---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_anycastip"
description: |-
  Manages a anycast IP resource within VKCS.
---

# vkcs_networking_anycastip



## Example Usage
### Anycast IP association with two octavia loadbalancers
```terraform
resource "vkcs_networking_anycastip" "anycastip" {
  name        = "app-tf-example"
  description = "app-tf-example"

  network_id = data.vkcs_networking_network.extnet.id
  associations = [
    {
      id   = vkcs_lb_loadbalancer.app1.vip_port_id
      type = "octavia"
    },
    {
      id   = vkcs_lb_loadbalancer.app2.vip_port_id
      type = "octavia"
    }
  ]
}
```

## Argument Reference
- `network_id` **required** *string* &rarr;  ID of the external network to choose ip for anycast IP from.

- `associations`  *set* &rarr;  List of port associations with anycast IP.
  - `id` **required** *string* &rarr;  ID of port / dc interface / octavia loadbalancer vip port.

  - `type` **required** *string* &rarr;  Type of association. Can be one of: port, dc_interface, octavia.


- `description` optional *string* &rarr;  Description of the anycast IP.

- `health_check` optional &rarr;  Health check settings.
  - `port` optional *number* &rarr;  Port for check to connect to.

  - `type` optional *string* &rarr;  Check type. Can be one of: TCP, ICMP.


- `name` optional *string* &rarr;  Name of the anycast IP.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the anycast IP.

- `ip_address` *string* &rarr;  Anycast IP address.

- `subnet_id` *string* &rarr;  Anycast IP subnet id.



## Import

Anycast IPs can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_anycastip.anycastip_1 bfbed405-dd89-41d9-aa97-6e335161146d
```
