---
layout: "vkcs"
page_title: "Setting up Data Platform products"
description: |-
  Step-by-step instruction for setting up Data Platform products in VKCS.
---

# Setting up Data Platform products

## Prerequisites

Before diving into the guide, ensure you meet the following prerequisites:

- **Configure Terraform and VKCS Provider** Make sure that you installed Terraform CLI and configured VKCS Provider. Follow [instructions](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs/guides/getting_started) if needed.
- **Understand Terraform Basics** Familiarize yourself with Terraform concepts like resource lifecycles, dependencies, and state management. [Terraform documentation](https://developer.hashicorp.com/terraform/docs) will help you understand the basic principles and key points.
- Make sure you have enough quotas to run examples of the guide.
- Have a look at `vkcs_dataplatform_cluster` arguments.


## Learn more about required products in VKCS

Look for required products in Data Platform section of [documentation](https://cloud.vk.com/docs/en). Review the products documentation to figure out configuration possibilities for the products.

Verify that necessary products are available in your project. Select Data Platform in Management Console and look for the neccesary products in a list of available products.


## Discover names and versions of required products

At the first you need to figure out the products' name and versions which are available in Terraform. Create `vkcs_dataplatform_products` data source and output for it which filters required products. In this guide we set up ClickHouse and Trino, so only this products are interesting.

```terraform
data "vkcs_dataplatform_products" "products" {}
```
```terraform
output "guide_products" {
  value = [
    for p in data.vkcs_dataplatform_products.products.products :
    p if strcontains(p.product_name, "click") || strcontains(p.product_name, "trino")
  ]
}
```

Here is the output value:
```
guide_products = [
  {
    "product_name" = "trino"
    "product_version" = "0.468.1"
  },
  {
    "product_name" = "trino"
    "product_version" = "0.468.0"
  },
  {
    "product_name" = "clickhouse"
    "product_version" = "25.3.0"
  },
  {
    "product_name" = "clickhouse"
    "product_version" = "24.3.0"
  },
]
```
Now you know technical names and versions of required projects to use them later.

**Note:** There are a lot of other products in the data source but avoid use those of them which are not available in your project. Because of they will not be supported for you.


## Discover data for pod_groups argument

`pod_groups` argument of `vkcs_dataplatform_cluster` resources defines parameters of nodes the product consists of. These parameters have an impact on resource capacity of product instance to handle user requests and data.

Lets discover what pod groups the products consist of. Create `vkcs_dataplatform_template` data source and output for names of pod groups. Here are examples for ClickHouse only but you could do the same for Trino by your own.

```terraform
data "vkcs_dataplatform_template" "clickhouse" {
  product_name    = "clickhouse"
  product_version = "25.3.0"
}
```
```terraform
output "clickhouse_pod_groups" {
  value = data.vkcs_dataplatform_template.clickhouse.pod_groups[*].name
}
```

Here is the output value:
```
clickhouse_pod_groups = tolist([
  "clickhouseKeeper",
  "clickhouse",
])
```

Data Platform has default configuration for pod groups. So you do not have to configure all of them to set up a product. However you should at least check default configuration. Lets see configuration of main ClickHouse nodes.

```terraform
output "clickhouse_pod_group" {
  value = one([
    for g in data.vkcs_dataplatform_template.clickhouse.pod_groups :
    g if g.name == "clickhouse"
  ])
}
```

Here is the output value:
```
clickhouse_pod_group = {
  "count" = 3
  "name" = "clickhouse"
  "resource" = {
    "cpu_margin" = 1
    "cpu_request" = "2.0"
    "ram_margin" = 1
    "ram_request" = "4.0"
  }
  "volumes" = tomap({
    "data" = {
      "count" = 1
      "storage" = "10"
      "storage_class_name" = "ceph-ssd"
    }
  })
}
```
See also `pod_groups` argument structure and description in documentation of `vkcs_dataplatform_cluster` resource.

**NOTE:** If you do not want to change default configuration of a pod group do not omit it but include in your manifest with default values. Otherwise you cannot scale the pod group later.

Now your know how to discover pod groups of products, their available and default configurations.

## Discover product configuration settings

`configs.settings` argument of `vkcs_dataplatform_cluster` resource lets you configure product software itself witch is under the hood of Data Platform product (i.e. ClickHouse, Trino, etc).

Lets see available settings and their default values. Create `vkcs_dataplatform_product` data source and outputs for default values of settings. Here are examples for ClickHouse only but you could do the same for Trino by your own.

```terraform
data "vkcs_dataplatform_product" "clickhouse" {
  product_name    = "clickhouse"
  product_version = "25.3.0"
}
```
```terraform
output "clickhouse_settings_all" {
  value = {
    for s in data.vkcs_dataplatform_product.clickhouse.configs.settings :
    s.alias => s.default_value
    # Settings which first part of name ends with "_" are not available for tuning
    if !endswith(split(".", s.alias)[0], "_")
  }
}
```

Here is partial output value:
```
clickhouse_settings_all = {
  "clickhouse.background_buffer_flush_schedule_pool_size" = "16"
  "clickhouse.background_common_pool_size" = "8"
  "clickhouse.background_distributed_schedule_pool_size" = "16"
...
}
```

For example lets dive to `clickhouse.background_common_pool_size` setting metadata. Create output for it.
```terraform
output "clickhouse_settings_background_common_pool_size" {
  value = one([
    for s in data.vkcs_dataplatform_product.clickhouse.configs.settings :
    s if s.alias == "clickhouse.background_common_pool_size"
  ])
}
```

Here is the metadata for the setting:
```
clickhouse_settings_background_common_pool_size = {
  "alias" = "clickhouse.background_common_pool_size"
  "default_value" = "8"
  "is_require" = false
  "is_sensitive" = false
  "regexp" = ""
  "string_variation" = tolist([])
}
```
The meaning of metadata attributes is pretty clear.

Now you know how to discover configuration settings of products, their default values and values restrictions.


## Discover product user roles

In `configs.users` argument of `vkcs_dataplatform_cluster` resource you specify user accounts to work with the product. Some Data Platform products allow role assignment for the accounts.

Lets explore available user roles. Create output of `vkcs_dataplatform_product` with role names. Here are examples for ClickHouse only but you could do the same for Trino by your own.

```terraform
output "clickhouse_user_roles" {
  value = data.vkcs_dataplatform_product.clickhouse.configs.user_roles[*].name
}
```

Here is the output value:
```
clickhouse_user_roles = tolist([
  "dbOwner",
  "readWrite",
  "readOnly",
  "secOps",
])
```

Now you know how to discover available user roles of products.


## Discover product connections

Some Data Platform products can connect to other data sources. Connections are specified with `config.warehouses.connection` argument of `vkcs_dataplatform_cluster` resource. `plug` argument means target type of a connection.

To discover available connection target types for the product create output of `vkcs_dataplatform_product` with connection `plug` value. Also take into account `is_required` attribute. Here are examples for Trino since in the guide we want to connect it to ClickHouse but you could do the same for ClickHouse by your own.

```terraform
data "vkcs_dataplatform_product" "trino" {
  product_name    = "trino"
  product_version = "0.468.1"
}
```
```terraform
output "trino_connections" {
  value = {
    for c in data.vkcs_dataplatform_product.trino.configs.connections :
    c.plug => c.is_required
  }
}
```

Here is the output value:
```
trino_connections = {
  "clickhouse" = false
  "greenplum" = false
  "hive" = false
  "hive-metastore-ext" = false
  "hive-metastore-int" = false
  "iceberg-metastore-ext" = false
  "iceberg-metastore-int" = false
  "mariadb" = false
  "mysql" = false
  "postgresql" = false
  "redis" = false
}
```

Lets dive into connection attributes.

```terraform
output "trino_connection_clickhouse" {
  value = one([
    for c in data.vkcs_dataplatform_product.trino.configs.connections :
    c.settings
    if c.plug == "clickhouse"
  ])
}
```

Here is partial output value with metadata of `username` connection attribute:
```
trino_connection_clickhouse = tolist([
...
  {
    "alias" = "username"
    "default_value" = ""
    "is_require" = true
    "is_sensitive" = false
    "regexp" = "^[a-zA-Z_][a-zA-Z0-9_]{1,49}$"
    "string_variation" = tolist([])
  },
...
])
```
The meaning of metadata attributes is pretty clear.

Now you know how to discover connection types of products, available connection attributes and thier values restrictions.


## Discover crontabs for product maintenance

Some Data Platform products have additional crontab settings besides maintenance period. See `config.maintenance.crontabs` argument of `vkcs_dataplatform_cluster` resource. Create output for them.

```terraform
output "trino_crontabs" {
  value = data.vkcs_dataplatform_product.trino.configs.crontabs[*]
}
```

Here is partial output value with metadata of the first setting:
```
trino_crontabs = tolist([
  {
    "name" = "maintenance"
    "required" = true
    "settings" = tolist([
      {
        "alias" = "duration"
        "default_value" = "3600"
        "is_require" = true
        "is_sensitive" = false
        "regexp" = ""
        "string_variation" = tolist([])
      },
...
    ])
    "start" = "0 22 * * *"
  },
])
```
Here is some mess. `is_require` does not mean you must specify this setting. `default_value` will be used if you omit this setting.

Now you know how to discover additional crontabs and their settings.


## Create manifests

Now you are ready to fill all needed arguments for resources of the guide.

```terraform
resource "vkcs_dataplatform_cluster" "clickhouse" {
  name            = "clickhouse-tf-guide"
  description     = "ClickHouse example instance from Data Platform guide."
  product_name    = "clickhouse"
  product_version = "25.3.0"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"
  # Enable public access to simplify testing of the product.
  floating_ip_pool = "auto"

  pod_groups = [
    # Omit settings for clickhouseKeeper pod group to illustrate
    # how Data Platform handle this with default settings.
    # NOTE: If you omit settings for a pod group you cannot scale
    # the pod group later.
    # Increase ram_request and storage values for clickhouse pod group
    # against default settings.
    {
      count = 3
      name  = "clickhouse"
      resource = {
        cpu_request = "2.0"
        ram_request = "8.0"
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "150"
          count              = 1
        }
      }
    },
  ]
  configs = {
    settings = [
      # Increase value of the setting against default one.
      {
        alias = "clickhouse.background_common_pool_size"
        value = 10
      },
    ]
    users = [
      {
        username = "owner"
        password = random_password.clickhouse_owner.result
        role     = "dbOwner"
      },
      {
        username = "trino"
        password = random_password.clickhouse_trino.result
        role     = "readOnly"
      },
    ]
    warehouses = [
      # Define database name.
      {
        name = "clickhouse"
      },
    ]
    maintenance = {
      # Set start om maintenance the same as start of full backup.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
      start = "0 1 * * 0"
      backup = {
        full = {
          keep_count = 5
          start      = "0 1 * * 0"
        }
        incremental = {
          keep_count = 7
          start      = "0 1 * * 1-6"
        }
      }
    }
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  depends_on = [vkcs_networking_router_interface.db]
}
```
```terraform
locals {
  # Parse ClickHouse endpoints to find tcp one.
  clickhouse_tcp = one([
    for srv in vkcs_dataplatform_cluster.clickhouse.info.services :
    {
      host = regex(".*@([^:/]+):([0-9]+).*", srv.connection_string)[0]
      port = regex(".*@([^:/]+):([0-9]+).*", srv.connection_string)[1]
    }
    if srv.type == "connection_string" && endswith(srv.description, "tcp")
  ])
}
```
```terraform
resource "vkcs_dataplatform_cluster" "trino" {
  name            = "trino-tf-guide"
  description     = "Trino example instance from Data Platform guide."
  product_name    = "trino"
  product_version = "0.468.1"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

  # In order to create Trino in the same cluster stack as the ClickHouse.
  stack_id = vkcs_dataplatform_cluster.clickhouse.stack_id
  # This argument must be the same for all products in the same cluster stack.
  floating_ip_pool = "auto"

  pod_groups = [
    {
      name  = "coordinator"
      count = 1
      resource = {
        cpu_request = "2.0"
        ram_request = "4.0"
      }
    },
    {
      name  = "worker"
      count = 1
      resource = {
        cpu_request = "2.0"
        ram_request = "4.0"
      }
    }
  ]
  configs = {
    users = [
      {
        username = "example"
        password = random_password.trino_example.result
      }
    ]
    warehouses = [{
      # For some Data Platform product value of the `name` argument has no sense
      # but is fixed. So you must set exactly this value.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
      name = "trino"
      connections = [
        {
          name = "clickhouse"
          plug = "clickhouse"
          settings = [
            {
              alias = "hostname"
              value = "${local.clickhouse_tcp.host}:${local.clickhouse_tcp.port}"
            },
            {
              alias = "username"
              value = "trino"
            },
            {
              alias = "password"
              value = random_password.clickhouse_trino.result
            },
            {
              alias = "ssl"
              value = "false"
            },
            {
              alias = "db_name"
              value = "clickhouse"
            },
            {
              alias = "catalog"
              value = "clickhouse"
            },
          ]
        },
      ]
    }]
    maintenance = {
      start = "0 22 * * *"
      crontabs = [
        {
          name  = "maintenance"
          start = "0 19 * * *"
        }
      ]
    }
  }
}
```
