---
layout: "vkcs"
page_title: "Building VPN Tunnel between VK Cloud private subnets"
description: |-
  VPN Tunnel within VKCS.
---

## VPN tunnel

VPN allows you to organize a tunnel between one or more VK Cloud private subnets and the client local network.
See https://cloud.vk.com/docs/en/networks/vnet/concepts/vpn for more info.

The tunnel is organized between private router and local VPN service. So you need:
- connect your private router to Internet;
- connect required private subnets to the router.

Building VPN tunnel depends on VKCS SDN you are using for networking.
Consult to `vkcs_networking_sdn` datasource to figure out available SDNs in you project.
If both `neutron` and `sprut` are available look at `sdn` attribute of network objects you are going to plug to VPN tunnel.

## Configure VPN tunnel

In this guide, we configure a vpn tunnel between Neutron network and Sprut network.
Complete example files can be viewed in the `examples/vpnaas` folder on GitHub project.
You need both SDNs enabled in your project to try the example.

### Neutron SDN:
If Neutron SDN is the only SDN available in a project then a standard router (`vkcs_networking_router`)
must be used as a private router to build VPN tunnel for Neutron networks.
Connect them to the router and specify router ID in `router_id` argument of `vkcs_vpnaas_service` resource.

Only `subnet` type is supported for `type` argument of private `vkcs_vpnaas_endpoint_group` resource.
Neutron VPNaaS can plug several subnets to one VPN tunnel but a way it is used for that could be not supported by certain client VPN hardware.
So check documentation and use the single subnet in this case.

If Sprut SDN is also available in the project then Sprut's VPN can used to build VPN tunnel for Neutron networks as well.

In this guide we build pure Neutron SDN VPN tunnel. In `examples/vpnaas/neutron_base.tf` file we prepare Neutron networks and router.
Next in `examples/vpnaas/neutron_main.tf` we build VPN tunnel.

```terraform
resource "vkcs_vpnaas_ike_policy" "neutron" {
  name        = "neutron-ike-tf-example"
  ike_version = "v2"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm          = local.auth_algorithm
  encryption_algorithm    = local.encryption_algorithm
  phase1_negotiation_mode = local.phase1_negotiation_mode
  sdn                     = local.sdn_neutron
}

resource "vkcs_vpnaas_ipsec_policy" "neutron" {
  name = "neutron-ipsec-tf-example"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm       = local.auth_algorithm
  encryption_algorithm = local.encryption_algorithm
  pfs                  = local.pfs
  sdn                  = local.sdn_neutron
}

# Local endpoint
resource "vkcs_vpnaas_endpoint_group" "local_neutron" {
  name = "neutron-local-endpoint-tf-example"
  type = "subnet"
  endpoints = [
    vkcs_networking_subnet.neutron.id
  ]
  sdn = local.sdn_neutron

  # Wait for subnet is connected to the router
  depends_on = [
    vkcs_networking_router_interface.neutron
  ]
}

# Remote endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_neutron" {
  name = "neutron-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.sprut_cidr
  ]
  sdn = local.sdn_neutron
}

resource "vkcs_vpnaas_service" "neutron" {
  name      = "neutron-vpn-tf-example"
  router_id = vkcs_networking_router.neutron.id
  sdn       = local.sdn_neutron
}

resource "vkcs_vpnaas_site_connection" "connection_neutron" {
  name              = "connection-neutron-tf-example"
  ikepolicy_id      = vkcs_vpnaas_ike_policy.neutron.id
  ipsecpolicy_id    = vkcs_vpnaas_ipsec_policy.neutron.id
  vpnservice_id     = vkcs_vpnaas_service.neutron.id
  psk               = local.psk_key
  peer_address      = vkcs_dc_interface.internet_sprut.ip_address
  peer_id           = vkcs_dc_interface.internet_sprut.ip_address
  local_ep_group_id = vkcs_vpnaas_endpoint_group.local_neutron.id
  peer_ep_group_id  = vkcs_vpnaas_endpoint_group.remote_neutron.id
  sdn               = local.sdn_neutron

  dpd {
    action   = "restart"
    timeout  = 60
    interval = 30
  }
}

# Static route in subnet to local network is very recommended.
# See https://cloud.vk.com/docs/networks/vnet/how-to-guides/vpn-tunnel#4_add_static_routes
# Find IP address of SNAT port of the router
data "vkcs_networking_port" "snat_port_neutron" {
  network_id   = vkcs_networking_network.neutron.id
  device_owner = "network:router_centralized_snat"

  # Wait for subnet is connected to the router
  depends_on = [
    vkcs_networking_router_interface.neutron
  ]
}

resource "vkcs_networking_subnet_route" "static_neutron" {
  subnet_id        = vkcs_networking_subnet.neutron.id
  destination_cidr = vkcs_networking_subnet.sprut.cidr
  next_hop         = data.vkcs_networking_port.snat_port_neutron.all_fixed_ips[0]
}
```

### Sprut SDN:
An advanced router (DC router, `vkcs_dc_router`) is required as a private router to build VPN tunnel for Sprut networks.
Supply network connectivity to the router for them and specify router ID in `router_id` argument of `vkcs_vpnaas_service` resource.

Only `cidr` type is supported for `type` argument of private `vkcs_vpnaas_endpoint_group` resource.
This makes VPN tunnel settings more flexible.

In the guide's `examples/vpnaas/sprut_base.tf` file we prepare advanced router and Sprut networks connected to it.
Next in `examples/vpnaas/sprut_main.tf` we build VPN tunnel.

```terraform
resource "vkcs_vpnaas_ike_policy" "sprut" {
  name        = "sprut-ike-tf-example"
  ike_version = "v2"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm          = local.auth_algorithm
  encryption_algorithm    = local.encryption_algorithm
  phase1_negotiation_mode = local.phase1_negotiation_mode
  sdn                     = local.sdn_sprut
}

resource "vkcs_vpnaas_ipsec_policy" "sprut" {
  name = "sprut-ipsec-tf-example"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm       = local.auth_algorithm
  encryption_algorithm = local.encryption_algorithm
  pfs                  = local.pfs
  sdn                  = local.sdn_sprut
}

# Local sprut endpoint
resource "vkcs_vpnaas_endpoint_group" "local_sprut" {
  name = "sprut-local-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.sprut_cidr
  ]
  sdn = local.sdn_sprut
}

# Remote sprut endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_sprut" {
  name = "sprut-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.neutron_cidr
  ]
  sdn = local.sdn_sprut
}

resource "vkcs_vpnaas_service" "sprut" {
  name      = "sprut-vpn-tf-example"
  router_id = vkcs_dc_router.sprut.id
  sdn       = local.sdn_sprut

  depends_on = [
    vkcs_dc_interface.internet_sprut
  ]
}

resource "vkcs_vpnaas_site_connection" "connection_sprut" {
  name              = "connection-sprut-tf-example"
  ikepolicy_id      = vkcs_vpnaas_ike_policy.sprut.id
  ipsecpolicy_id    = vkcs_vpnaas_ipsec_policy.sprut.id
  vpnservice_id     = vkcs_vpnaas_service.sprut.id
  psk               = local.psk_key
  peer_address      = vkcs_networking_router.neutron.external_fixed_ips[0].ip_address
  peer_id           = vkcs_networking_router.neutron.external_fixed_ips[0].ip_address
  local_ep_group_id = vkcs_vpnaas_endpoint_group.local_sprut.id
  peer_ep_group_id  = vkcs_vpnaas_endpoint_group.remote_sprut.id
  sdn               = local.sdn_sprut

  dpd {
    action   = "restart"
    timeout  = 60
    interval = 30
  }
}

# We have to specify a path to DC router for VPN traffic
# since it is not a gateway in the network.
resource "vkcs_networking_subnet_route" "static_sprut" {
  subnet_id        = vkcs_networking_subnet.sprut.id
  destination_cidr = local.neutron_cidr
  next_hop         = vkcs_dc_interface.subnet_sprut.ip_address
}
```

## Verify tunnel
In the guide's example we start two VMs in Neutron and Sprut networks and check connectivity between them.
See `examples/vpnaas/neutron_vm.tf`, `examples/vpnaas/sprut_vm.tf` if you are interested.
