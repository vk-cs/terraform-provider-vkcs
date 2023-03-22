---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog

#### v0.1.15 (unreleased)
- Add deprecation warning to "security_group_ids" field of resource_vkcs_lb_loadbalancer
- Add "instances" computed field to resource_vkcs_db_cluster_with_shards.shard
- Add "loadbalancer_id" computed field to resource_vkcs_db_cluster
- Fix error of not expecting "retyping" status when modifying resource_vkcs_blockstorage_volume.volume_type
- Add vkcs_lb_loadbalancer datasource
- Make "export_location_path" field of data_source_vkcs_sharedfilesystem_share computed
- Add "export_location_path" computed field to resource_vkcs_sharedfilesystem_share
- Fix error of ignoring restore_point field in resource_vkcs_db_instance and resource_vkcs_db_cluster

#### v0.1.14
- Add "instances" computed field to resource_vkcs_kubernetes_cluster
- Add ability to control which cluster instances should remain after shrinking cluster via "shrink_options" field of resource_vkcs_kubernetes_cluster
- Add "ip" computed field to resource_vkcs_db_instance
- Allow creation of resource_vkcs_images_image in clouds without s3 support
- Added description for cluster_node_volume_type labels field in resource_vkcs_kubernetes_cluster

#### v0.1.13
- Updated description for labels field in resource_vkcs_kubernetes_cluster 
- Added conflicts_with property to remote_group_id and remote_ip_prefix fields of resource_vkcs_networking_secgroup_rule
- Added deprecation warning to ethertype field of resource_vkcs_networking_secgroup_rule
- Removed dns_name and dns_domain attributes from resource_vkcs_networking_floatingip
- Removed ip_version parameter from networking resources
- Removed ipv6 related parameters from vkcs_compute_instance resource and data_source, from networking resources and data_sources
- Removed sort_key and sort_direction arguments from data_source_vkcs_images_image
- Removed parent_region_id argument from data_source_vkcs_regions
- Removed cacert_file, cert, insecure, key fields from provider config

#### v0.1.12
- Removed image_id attribute from resource_vkcs_images_image and hidden attribute from vkcs_images_image resource and data_source
- Removed tags attribute from vkcs_networking_floatingip resource and data_source
- Fix getting OS_USERNAME, OS_REGION_NAME and OS_USER_DOMAIN_NAME from environment

#### v0.1.11
- Fixed error of creating resource_vkcs_blockstorage_volume with high-iops volume type

#### v0.1.10
- Minor updates for resource_vkcs_kubernetes_cluster, resource_vkcs_kubernetes_node_group,
- data_source_vkcs_kubernetes_cluster and examples

#### v0.1.9
- Fix re-creation of db user after re-creation of db instance 

#### v0.1.8
- Minor updates for resource_vkcs_kubernetes_cluster, resource_vkcs_kubernetes_node_group, 
- data_source_vkcs_kubernetes_cluster and examples

#### v0.1.7
- Add vkcs_db_config_group resource and datasource

#### v0.1.6
- Fixed error handling of creating root user for resource_vkcs_db_instance
- Minor updates to resource_vkcs_db_user and resource_vkcs_db_database documentation

#### v0.1.5
- Add dns_domain field to vkcs_kubernetes_cluster

#### v0.1.4
- Add vkcs_db_backup resource and datasource
- Add restore_point field to vkcs_db_instance, vkcs_db_cluster and vkcs_db_cluster_with_shards resources
- Add backup_schedule field to vkcs_db_instance resource and data_source and to vkcs_db_cluster resource

#### v0.1.3
- Add max_node_unavailable option for node group

#### v0.1.2
- Fix creation of db databases and users

#### v0.1.1
- Fix typo in data-sources docs location
- Fix database instances import: databases & users

#### v0.1.0
- Initial release
