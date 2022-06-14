---
layout: "vkcs"
page_title: "VKCS provider changelog"
description: |-
  The VKCS provider's changelog.
---

# VKCS Provider's changelog

#### v.0.1.5
- Add dns_domain field to vkcs_kubernetes_cluster

#### v.0.1.4
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
