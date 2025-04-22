---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog

#### v0.9.4 (unreleased)

#### v0.9.3
- Add secure_key option to cdn resource
- Add gzip_compression option to cdn resource
- Add availability_zones and vrrp_port_id fields to vkcs_db_cluster resource
- Add ability go get admin password of compute instance with windows
- Support multi-AZ PostgreSQL clusters
- Fix filtering in the vkcs_cdn_shielding_pop data source
- Mark api_lb_vip field of vkcs_kubernetes_cluster resource as read-only
- Deprecate fixed_ip_v4 field of vkcs_db_instance resource

#### v0.9.2
- Reduce the minimum number of DHCP ports to wait for a subnet to be ready from 2 to 1
- Remove retries for bad request error when attaching a volume to a Compute instance

#### v0.9.1
- Fix panic on VKCS Kubernetes API error in vkcs_kubernetes_security_policy
- Update Let's Encrypt certificate issuance process in vkcs_cdn_resource

#### v0.9.0
- Add vkcs_cloud_monitoring resource
- Add cloud_monitoring argument to vkcs_compute_instance
- Add support for CDN service
- Fix an issue with marking subnet as ready while DHCP service is not configured yet
- Fix panic on empty taint to vkcs_kubernetes_node_group resource
- Fix order-induced changes in the plan for "databases" of the vkcs_db_user resource
- Fix panic on VKCS Kubernetes API error in vkcs_kubernetes_clustertemplates data source
- Fix an error when changing several Kubernetes cluster node groups in parallel
- Fix an error when reading the state of a compute instance with the build status
- Fix unexpected state error when the database cluster was performing backups after creation
- Add the ability to specify the id field instead of \<resource\>_id in data sources.
- Deprecate \<resource\>_id field of data sources.
- Deprecate ingress_floating_ip field of vkcs_kubernetes_cluster resource and data source
- Deprecate keypair field of vkcs_kubernetes_cluster resource and data source

#### v0.8.4
- Add ability to import vkcs_kubernetes_security_policy into the state 
- Fix forced re-creation of the vkcs_kubernetes_cluster due to redundant planned changes for availability_zones argument
- Fix order-induced changes in the plan for "availability_zones" of the vkcs_kubernetes_node_group resource

#### v0.8.3
- Add ability to choose type of kubernetes cluster, standard or regional
- Add cluster_type and availability_zones argument to vkcs_kubernetes_cluster
- Fix unexpected state error when updating volume_type of attached vkcs_blockstorage_volume
- Add a new pending status for the vkcs_blockstorage_volume extending
- Add all_metadata computed field to vkcs_blockstorage_volume resource
- Suppress out-of-scope plan changes in vkcs_blockstorage_volume metadata
- Increased the timeout for create and update operations for vkcs_sharedfilesystem_share

#### v0.8.2
- Add retry on duplicate IpamAllocation error when creating a networking router
- Stabilize vkcs_kubernetes_nodegroup creation

#### v0.8.1
- Add "k8s_config" attribute to vkcs_kubernetes_cluster resource
- Fix order-induced changes in the plan for "allowed_cidrs" of the vkcs_lb_listener resource
- Fix an error when filtering by pool name in the vkcs_networking_floatingip data source in SDN Sprut

#### v0.8.0
- Add vkcs_dc_conntrack_helper resource
- Add vkcs_dc_ip_port_forwarding resource
- Add full_security_groups_control argument to vkcs_networking_port resource.
- Deprecate no_security_groups argument of vkcs_networking_port resource.
- Added waiting for roles to be assigned to instance when creating a db cluster. 

#### v0.7.4
- Add vkcs_networking_sdn data source for getting a list of available SDNs.
- Add traffic_selector_ep_merge argument to vkcs_vpnaas_site_connection resource.
- Add external_fixed_ips computed field to vkcs_networking_router resource and data source.
- Update the Compute API client micro version to 2.42
- Fix searching for a subnet with a filter by "tags" in the vkcs_networking_subnet data source
- Fix unnecessary resource recreation due to shelve/unshelve operations on vkcs_compute_instance

#### v0.7.3
- Add security_group_ids argument to vkcs_compute_instance resource
- Add ability to filter by extra_specs attribute in vkcs_compute_flavor data source.
- Enhancement searching the closest appropriate flavor by min_ram argument in vkcs_compute_flavor data source.

#### v0.7.2
- Allow 'ip_address' parameter for vkcs_dc_interface resource creation

#### v0.7.1
- Remove reading of security groups on VIP port of load balancer
- Fix error on reading vkcs_networking_floatingip_associate resource
- Omit sending SDN header for operations on existing resources

#### v0.7.0
- Add vkcs_kubernetes_security_policy resource
- Add vkcs_kubernetes_security_policy_template and vkcs_kubernetes_security_policy_templates data sources
- Add 'sync_security_policy' parameter for vkcs_kubernetes_cluster resource and data source
- Allow getting deprecated images by vkcs_images_image data source
- Add vkcs_mlplatform_spark_k8s and vkcs_mlplatform_k8s_registry resources

#### v0.6.1
- Allow usage of Postgres Pro Enterprise and Postgres Pro Enterprise 1C datastore for database clusters and instance with replicas
- Make 'private_key' attribute of resource vkcs_compute_keypair sensitive

#### v0.6.0
- Add mlplatform jupyterhub, mlflow, mlflow_deploy resources
- Use chart name as Kubernetes addon's default name

#### v0.5.5
- Add force_deletion option to vkcs_db_database resource

#### v0.5.4
- Fix creating vkcs_db_backup for clusters
- Fix error when trying to read vkcs_db_cluster without instances
- Fix provider crash if there was a problem with authentication

#### v0.5.3
- Add 'sdn' parameter for vpnaas resources
- Add "skip_deletion" option to control the deletion strategy for Databases user resource
- Add retrying for 5xx http errors
- Fix error of inconsistent result on apply when creating a Kubernetes addon
- Fix error with creating large number of resources at once

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
