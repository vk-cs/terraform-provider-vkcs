---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_cluster"
description: |-
  Manages a dataplatform cluster
---

# vkcs_dataplatform_cluster



## Example Usage

# Opensearch
```terraform
resource "vkcs_dataplatform_cluster" "basic-opensearch" {
  cluster_template_id = "17230566-bfaa-492c-b22c-33b500a17155"
  name                = "tf-basic-opensearch"
  network_id          = "d2fad739-1b10-4dc8-9b2c-c246d7a7cc69"
  subnet_id           = "3a744943-fcc1-4a85-a96b-3dc4fff71885"
  product_name        = "opensearch"
  product_version     = "2.15.0"

  availability_zone = "UD2"
  configs = {
    users = [
      {
        user     = "user"
        password = "somepa55word!"
      }
    ]
    maintenance = {
      start = "0 22 * * *"
    }
    warehouses = [
      {
        name = "db_customer"
        users = [
          "user"
        ]
      }
    ]
  }
  pod_groups = [
    {
      count = 3
      resource = {
        cpu_request = "0.5"
        ram_request = "2"
      }
      volumes = [
        {
          type               = "data"
          count              = 3
          storage            = "30"
          storage_class_name = "ceph-hdd"
        }
      ]
      pod_group_template_id = "75cd99d3-7089-4372-8353-3ce675f55284"
    },
    {
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      pod_group_template_id = "1da36325-1bc8-43eb-ab52-113332001cff"
    }
  ]

}
```

# Trino
```terraform
resource "vkcs_dataplatform_cluster" "basic-trino" {
  cluster_template_id = "17230566-bfaa-492c-b22c-33b500a17155"
  name                = "tf-basic-trino"
  network_id          = "d2fad739-1b10-4dc8-9b2c-c246d7a7cc69"
  subnet_id           = "3a744943-fcc1-4a85-a96b-3dc4fff71885"
  product_name        = "trino"
  product_version     = "0.468.0"

  availability_zone = "UD2"
  configs = {
    users = [
      {
        user     = "user"
        password = "somepa55word!"
      }
    ]
    maintenance = {
      start = "0 22 * * *"
    }
    warehouses = [
      {
        name = "db_customer"
        users = [
          "user"
        ]
      }
    ]
  }
  pod_groups = [
    {
      count = 1
      resource = {
        cpu_request = "4"
        ram_request = "6"
      }
      pod_group_template_id = "75cd99d3-7089-4372-8353-3ce675f55284"
    },
    {
      count = 0
      resource = {
        cpu_request = "4"
        ram_request = "6"
      }
      pod_group_template_id = "1da36325-1bc8-43eb-ab52-113332001cff"
    }
  ]
}
```

## Argument Reference
- `cluster_template_id` **required** *string* &rarr;  ID of the cluster template.

- `name` **required** *string* &rarr;  Name of the cluster.

- `network_id` **required** *string* &rarr;  ID of the cluster network.

- `product_name` **required** *string* &rarr;  Name of the product.

- `product_version` **required** *string* &rarr;  Version of the product.

- `availability_zone` optional *string* &rarr;  Availability zone to create cluster in.

- `configs` optional &rarr;  Product configuration.
  - `maintenance` ***required*** &rarr;  Maintenance settings.
    - `backup` optional &rarr;  Backup settings.
      - `differential` optional &rarr;  Differential backup settings.
        - `start` **required** *string* &rarr;  Differential backup schedule.

        - `keep_count` optional *number* &rarr;  Differential backup keep count.

        - `keep_time` optional *number* &rarr;  Differential backup keep time.

        - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.


      - `full` optional &rarr;  Full backup settings.
        - `start` **required** *string* &rarr;  Full backup schedule.

        - `keep_count` optional *number* &rarr;  Full backup keep count.

        - `keep_time` optional *number* &rarr;  Full backup keep time.

        - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.


      - `incremental` optional &rarr;  Incremental backup settings.
        - `start` **required** *string* &rarr;  Incremental backup schedule.

        - `keep_count` optional *number* &rarr;  Incremental backup keep count.

        - `keep_time` optional *number* &rarr;  Incremental backup keep time.

        - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.



    - `cron_tabs`  *set* &rarr;  Cron tabs settings.
      - `name` **required** *string* &rarr;  Cron tab name.

      - `settings`  *set* &rarr;  Additional cron settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `start` optional *string* &rarr;  Cron tab schedule.

      - `required` read-only *boolean* &rarr;  Whether cron tab is required.


    - `start` optional *string* &rarr;  Maintenance cron schedule.


  - `features` optional &rarr;  Product features.
    - `volume_autoresize` optional &rarr;  Volume autoresize options.
      - `data` optional &rarr;  Data volume options.
        - `enabled` optional *boolean* &rarr;  Enables option.

        - `max_scale_size` optional *number* &rarr;  Maximum scale size.

        - `scale_step_size` optional *number* &rarr;  Scale step size.

        - `size_scale_threshold` optional *number* &rarr;  Size scale threshold.


      - `wal` optional &rarr;  Wal volume options.
        - `enabled` optional *boolean* &rarr;  Enables option.

        - `max_scale_size` optional *number* &rarr;  Maximum scale size.

        - `scale_step_size` optional *number* &rarr;  Scale step size.

        - `size_scale_threshold` optional *number* &rarr;  Size scale threshold.




  - `settings`  *set* &rarr;  Additional common settings.
    - `alias` **required** *string* &rarr;  Setting alias.

    - `value` **required** *string* &rarr;  Setting value.


  - `users`  *set* &rarr;  Users settings.
    - `password` **required** *string* &rarr;  Password.

    - `user` **required** *string* &rarr;  Username.

    - `access` optional &rarr;  Access settings.
      - `settings`  *set* &rarr;  Access users settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `id` read-only *string* &rarr;  Access ID.


    - `created_at` read-only *string* &rarr;  User creation timestamp.

    - `user_id` read-only *string* &rarr;  User ID.


  - `warehouses`  *set* &rarr;  Warehouses settings.
    - `connections`  *set* &rarr;  Warehouse connections.
      - `name` **required** *string* &rarr;  Connection name.

      - `plug` **required** *string* &rarr;  Connection plug.

      - `settings`  *set* &rarr;  Additional warehouse settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `created_at` read-only *string* &rarr;  Connection creation timestamp.

      - `id` read-only *string* &rarr;  Connection ID.


    - `extensions`  *set* &rarr;  Warehouse extensions.
      - `type` **required** *string* &rarr;  Extension type.

      - `settings`  *set* &rarr;  Additional extension settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `version` optional *string* &rarr;  Extension version.

      - `created_at` read-only *string* &rarr;  Extension creation timestamp.

      - `id` read-only *string* &rarr;  Extension ID.


    - `name` optional *string* &rarr;  Warehouse name.

    - `users` optional *set of* *string* &rarr;  Warehouse users.

    - `id` read-only *string* &rarr;  Warehouse ID.



- `description` optional *string* &rarr;  Cluster description.

- `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID.

- `multi_az` optional *boolean* &rarr;  Enables multi az support.

- `pod_groups`  *set* &rarr;  Cluster pod groups.
  - `pod_group_template_id` **required** *string* &rarr;  Pod group template ID.

  - `count` optional *number* &rarr;  Pod count.

  - `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID.

  - `node_processes` optional *set of* *string* &rarr;  Node processes.

  - `resource` optional &rarr;  Resource request settings.
    - `cpu_request` optional *string* &rarr;  CPU request settings.

    - `ram_request` optional *string* &rarr;  RAM request settings.

    - `cpu_limit` read-only *string* &rarr;  CPU limit.

    - `ram_limit` read-only *string* &rarr;  RAM limit settings.


  - `volumes`  *set* &rarr;  Volumes settings.
    - `count` **required** *number* &rarr;  Volume count.

    - `storage` **required** *string* &rarr;  Storage size.

    - `storage_class_name` **required** *string* &rarr;  Storage class name.

    - `type` **required** *string* &rarr;  Volume type.


  - `id` read-only *string* &rarr;  Pod group ID.


- `region` optional *string* &rarr;  The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `stack_id` optional *string* &rarr;  ID of the cluster stack.

- `subnet_id` optional *string* &rarr;  ID of the cluster subnet.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Cluster creation timestamp.

- `id` *number* &rarr;  ID of the cluster.


