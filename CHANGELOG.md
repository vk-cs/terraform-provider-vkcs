---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog

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
