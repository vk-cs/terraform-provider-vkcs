---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog
#### v0.5.3
- Add retrying for 5xx http errors
- Fix error of inconsistent result on apply when creating a Kubernetes addon
- Add "skip_deletion" option to control the deletion strategy for Databases user resource
- Fix error with creating large number of resources at once
- Add 'sdn' parameter for some network resources (endpoint_group, service and site-connection)

#### v0.5.2
- Fixed issue when vkcs_networking_router_route could not be read

#### v0.5.1
- Fixed issue of root password not working for vkcs_db_cluster and vkcs_db_cluster_with_shards resources
- Remove default sdn value from firewall and networking resources and data sources
- Support updating Kubernetes cluster labels without recreating it
- Add "all_labels" attribute to Kubernetes cluster
- Support updating Kubernetes node group's flavor
- Wait for volume to finish detaching when deleting resource vkcs_compute_volume_attach

#### v0.5.0
- Add Direct Connect resources and data sources
- Fix Kubernetes addon state refresh panic on error
- Fix error of unexpected Kubernetes cluster status after creating/deleting an addon
- Fix panic on Kubernetes cluster "status" attribute change 

#### v0.4.2
- Add data source for images
- Add "default" attribute to Images data sources
- Add options to unzip and decompress downloaded image in Images image resource
- Fix leaving a damaged file in cache when an error ocurred during downloading an image from source
- Fix schedule.time format of vkcs_backup_plan datasource

#### v0.4.1
- Fix crash on errored vkcs_backup_plan creation
- Fix overriding userDomainName with defaults on devstack environment

#### v0.4.0
- Add Cloud Backup resource and datasources
- Add restart_confirmed vendor option to DB clusters and instance
- Allow set vkcs_images_image properties with os_ prefix
- Fix creating/updating DB clusters with specified configuration_id
- Fix error of unexpected "restart_required" status for DB instance

#### v0.3.0
- Provide support for Kubernetes cluster addons
- Mark resources created with error as "tainted" which allows their deleting via terraform
- Fix duplicate entry error on recreating Kubernetes cluster
- Fix error when creating DB instance with specified configuration_id
- Prevent panic when reading whether root user is enabled for DB instance

#### v0.2.2
- Remove "optional" property from ip attribute of DB instance
- Add validation and default value for block_device.boot_index of Compute instance resource
- Structure the documentation by grouping it by service
- Fix error of concurrent file read/write in Image resource
- Fix error when creating DB instance replica with specified configuration_id
- Fix errors related to unexpected new statuses for DB instance and DB backup resources

#### v0.2.1
- Support shards/instances/volumes resizing in DB cluster with shards
- Format the "ip" attribute of the Public DNS AAAA record to keep plan and after-apply state in sync
- Fix reading state of secgroup rule with specified ethertype
- Fix reading state of DB cluster with shards

#### v0.2.0
- Provide support for Public DNS service
- Add datasources for DB datastores
- Add datasources for DB datastore capabilities and configuration parameters
- Support cloud monitoring for DB instances/clusters
- Add security_group argument to DB instance/cluster resources
- Fix searching of shared network with datasource

#### v0.1.16
- Add config option to run against clouds with old cloud containers API

#### v0.1.15
- Fix error of incorrect setting a value of resource_vkcs_compute_instance.access_ip_v4 to an empty string
- Add "subnet_id" argument to resource_vkcs_db_instance.network, resource_vkcs_db_cluster.network, and resource_vkcs_db_cluster_with_shards.shard.network
- Add a warning that is thrown when arguments "network.fixed_ip_v4" and "replica_of" of resource_vkcs_db_instance are set simultaneously
- Deprecate "port" argument of resource_vkcs_db_instance.network, resource_vkcs_db_cluster.network, and resource_vkcs_db_cluster_with_shards.shard.network
- Fix error of not resolving referenced network resources that were created with sdn = "sprut"
- Deprecate "security_group_ids" argument of resource_vkcs_lb_loadbalancer
- Add "instances" computed attribute to resource_vkcs_db_cluster_with_shards.shard
- Add "loadbalancer_id" computed attribute to resource_vkcs_db_cluster
- Fix error of not expecting "retyping" status when modifying resource_vkcs_blockstorage_volume.volume_type
- Add vkcs_lb_loadbalancer datasource
- Make "export_location_path" attribute of data_source_vkcs_sharedfilesystem_share computed
- Add "export_location_path" computed attribute to resource_vkcs_sharedfilesystem_share
- Fix error of ignoring "restore_point" argument in resource_vkcs_db_instance and resource_vkcs_db_cluster

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
