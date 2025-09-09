package dataplatform_test

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"time"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"text/template"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/randutil"
)

func TestAccDataPlatformClusterResource_update_big(t *testing.T) {
	s3BucketSuffix := randutil.RandomName(5)
	oldName := "tf-acc-spark"
	newName := "tf-acc-spark-new"
	oldDescription := "tf-acc-spark-desc"
	newDescription := "tf-acc-spark-new-desc"
	oldVersion := "spark-py-3.4.2:v3.5.1.2"
	newVersion := "spark-py-3.5.1:v3.5.1.2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDataPlatformClusterResourceUpdate, map[string]string{
					"TestAccDataPlatformClusterResourceBaseNetwork": testAccDataPlatformClusterResourceBaseNetwork,
					"TestAccDataPlatformClusterResourceBaseDB":      testAccDataPlatformClusterResourceBaseDB,
					"Name":         oldName,
					"Description":  oldDescription,
					"SparkVersion": oldVersion,
					"S3Bucket":     fmt.Sprintf("tfacc-s3-%s", s3BucketSuffix),
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "name", oldName),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "description", oldDescription),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.settings.0.alias", "sparkproxy.spark_version"),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.settings.0.value", oldVersion),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDataPlatformClusterResourceUpdate, map[string]string{
					"TestAccDataPlatformClusterResourceBaseNetwork": testAccDataPlatformClusterResourceBaseNetwork,
					"TestAccDataPlatformClusterResourceBaseDB":      testAccDataPlatformClusterResourceBaseDB,
					"Name":         newName,
					"Description":  newDescription,
					"SparkVersion": newVersion,
					"S3Bucket":     fmt.Sprintf("tfacc-s3-%s", s3BucketSuffix),
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_dataplatform_cluster.basic", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "name", newName),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "description", newDescription),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.settings.0.alias", "sparkproxy.spark_version"),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.settings.0.value", newVersion),
				),
			},
		},
	})
}

type DataplatformClusterUser struct {
	Username string
	Password string
	Role     string
}

func testRenderDataplatformClusterUsers(users []DataplatformClusterUser) string {
	usersTmpl := template.Must(template.New("users").Option("missingkey=error").Parse(testAccDataPlatformClusterResourceIcebergUsers))
	var buf bytes.Buffer

	_ = usersTmpl.Execute(&buf, users)
	return buf.String()
}

func TestAccDataPlatformClusterIcebergAddAndDeleteUser_big(t *testing.T) {
	oneUser := testRenderDataplatformClusterUsers([]DataplatformClusterUser{{Username: "vkdata", Password: "Test_p#ssword-12-3", Role: "dbOwner"}})
	twoUsers := testRenderDataplatformClusterUsers([]DataplatformClusterUser{{Username: "vkdata", Password: "Test_p#ssword-12-3", Role: "dbOwner"}, {Username: "vkdata1", Password: "Test_p#ssword-12-4", Role: "common"}})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPlatformClusterResourceBaseNetwork,
				Check: func(state *terraform.State) error {
					time.Sleep(30 * time.Second)
					return nil
				},
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDataPlatformClusterResourceIceberg, map[string]string{
					"TestAccDataPlatformClusterResourceBaseNetwork":  testAccDataPlatformClusterResourceBaseNetwork,
					"TestAccDataPlatformClusterResourceIcebergUsers": oneUser,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "name", "tf-basic-iceberg"),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "description", "tf-basic-iceberg-description"),
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.users.#", "1"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDataPlatformClusterResourceIceberg, map[string]string{
					"TestAccDataPlatformClusterResourceBaseNetwork":  testAccDataPlatformClusterResourceBaseNetwork,
					"TestAccDataPlatformClusterResourceIcebergUsers": twoUsers,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.users.#", "2"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDataPlatformClusterResourceIceberg, map[string]string{
					"TestAccDataPlatformClusterResourceBaseNetwork":  testAccDataPlatformClusterResourceBaseNetwork,
					"TestAccDataPlatformClusterResourceIcebergUsers": oneUser,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dataplatform_cluster.basic", "configs.users.#", "1"),
				),
			},
		},
	})
}

const testAccDataPlatformClusterResourceBaseNetwork = `
resource "vkcs_networking_network" "db" {
  name        = "db-tf-example"
  description = "Database network"
}

resource "vkcs_networking_subnet" "db" {
  name       = "db-tf-example"
  network_id = vkcs_networking_network.db.id
  cidr       = "192.168.166.0/24"
}

data "vkcs_networking_network" "extnet" {
  name = "internet"
}

resource "vkcs_networking_router" "router" {
  name = "router-tf-example"
  # Connect router to Internet
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "db" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.db.id
}
`

const testAccDataPlatformClusterResourceBaseDB = `
data "vkcs_compute_flavor" "db" {
  name = "Standard-2-8-50"
}

resource "vkcs_db_instance" "db_instance" {
  name = "db-instance-tf-example"

  availability_zone = "GZ1"

  datastore {
    type    = "postgresql"
    version = "16"
  }

  flavor_id           = data.vkcs_compute_flavor.db.id
  floating_ip_enabled = true

  size        = 10
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

resource "vkcs_db_database" "postgres_db" {
  name    = "testdb_1"
  dbms_id = vkcs_db_instance.db_instance.id

  vendor_options {
    force_deletion = true
  }
}

resource "vkcs_db_user" "postgres_user" {
  name     = "testuser"
  password = "Test_p#ssword-12-3"

  dbms_id = vkcs_db_instance.db_instance.id

  vendor_options {
    skip_deletion = true
  }

  databases = [vkcs_db_database.postgres_db.name]
}
`

const testAccDataPlatformClusterResourceUpdate = `
{{ .TestAccDataPlatformClusterResourceBaseNetwork }}
{{ .TestAccDataPlatformClusterResourceBaseDB }}

resource "vkcs_dataplatform_cluster" "basic" {
  name            = "{{ .Name }}"
  description     = "{{ .Description }}"
  network_id      = vkcs_networking_network.db.id
  subnet_id       = vkcs_networking_subnet.db.id
  product_name    = "spark"
  product_version = "3.5.1"

  availability_zone = "GZ1"
  configs = {
    settings = [
      {
        alias = "sparkproxy.spark_version"
        value = "{{ .SparkVersion }}"
      }
    ]
    maintenance = {
      start = "0 0 1 * *"
    }
    warehouses = [
      {
        name = "spark"
        connections = [
          {
            name = "s3_int"
            plug = "s3-int"
            settings = [
              {
                alias = "s3_bucket"
                value = "{{ .S3Bucket }}"
              },
              {
                alias = "s3_folder"
                value = "tfacctest-folder"
              }
            ]
          },
          {
            name = "postgres"
            plug = "postgresql"
            settings = [
              {
                alias = "db_name"
                value = vkcs_db_database.postgres_db.name
              },
              {
                alias = "hostname"
                value = "${vkcs_db_instance.db_instance.ip[0]}:5432"
              },
              {
                alias = "username"
                value = vkcs_db_user.postgres_user.name
              },
              {
                alias = "password"
                value = vkcs_db_user.postgres_user.password
              }
            ]
          }
        ]
      }
    ]
  }
  pod_groups = [
    {
      name  = "sparkconnect"
      count = 1
      resource = {
        cpu_request = "10"
        ram_request = "10"
      }
    },
    {
      name  = "sparkhistory"
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      volumes = {
        "data" = {
          storage_class_name = "ceph-ssd"
          storage            = "5"
          count              = 1
        }
      }
    }
  ]
}
`

const testAccDataPlatformClusterResourceIcebergUsers = `
	users = [{{range .}}
	  {
		username = "{{.Username}}"
		password = "{{.Password}}"
		role = "{{.Role}}"
	  },{{end}}
	]
`

const testAccDataPlatformClusterResourceIceberg = `
{{ .TestAccDataPlatformClusterResourceBaseNetwork }}

resource "vkcs_dataplatform_cluster" "basic" {
  name            = "tf-basic-iceberg"
  description     = "tf-basic-iceberg-description"
  network_id      = vkcs_networking_network.db.id
  subnet_id       = vkcs_networking_subnet.db.id
  product_name    = "iceberg-metastore"
  product_version = "17.2.0"

  availability_zone = "GZ1"
  configs = {
    maintenance = {
      start = "0 0 1 * *"
      backup = {
        full = {
          keep_time = 10
          start = "0 0 1 * *"
        }
      }
    }
    warehouses = [
      {
        name = "metastore"
      }
    ]
    {{ .TestAccDataPlatformClusterResourceIcebergUsers }} 
  }
  pod_groups = [
    {
      name  = "postgres"
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      volumes = {
        "data" = {
          storage_class_name = "ceph-ssd"
          storage            = "10"
          count              = 1
        }
        "wal" = {
          storage_class_name = "ceph-ssd"
          storage = "10"
          count = 1
        }
      }
    },
    {
      name = "bouncer"
      count = 1
      resource = {
        cpu_request = "0.2"
        ram_request = "1"
      }
    }
  ]
}
`
