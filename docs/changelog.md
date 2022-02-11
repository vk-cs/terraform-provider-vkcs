---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog

#### v0.5.8
- Removed attribute `ingress_floating_ip` from `mcs_kubernetes_cluster`. 

#### v0.5.7
- Forbade using name in `master_flavor` attribute in `mcs_kubernetes_cluster`.
- Forbade using name in `flavor_id` attribute in `mcs_kubernetes_nodegroup`.

#### v0.5.6
- Make `name` attribute of node group required.

#### v0.5.4
- Added `loadbalancer_subnet_id` attribute to cluster.

#### v0.5.0
- Added `availability_zones` attribute to cluster node group.
- Added `mcs_kubernetes_clustertemplates` data source.

#### v0.4.0
- Added `region` support for provider.
- Added `mcs_region` and `mcs_regions` data sources.

#### v0.3.4
- Removed field `node_count` for kubernetes cluster.

#### v0.3.3
- Added required field `availablity_zone` to kubernetes cluster.