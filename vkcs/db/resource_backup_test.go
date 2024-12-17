package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseBackup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseBackupBasic, map[string]string{"TestAccDatabaseInstanceBasic": acctest.AccTestRenderConfig(testAccDatabaseInstanceBasic)}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_db_backup.basic", "name", "basic"),
				),
			},
		},
	})
}

func TestAccDatabaseBackupCluster_basic_big(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseBackupClusterBasic, map[string]string{"TestAccDatabaseClusterBasic": acctest.AccTestRenderConfig(testAccDatabaseClusterBasic)}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_db_backup.basic", "name", "basic"),
				),
			},
		},
	})
}

const testAccDatabaseBackupBasic = `
{{.TestAccDatabaseInstanceBasic}}

resource "vkcs_db_backup" "basic" {
    name = "basic"
    dbms_id = vkcs_db_instance.basic.id
    description = "basic description"
}
`

const testAccDatabaseBackupClusterBasic = `
{{.TestAccDatabaseClusterBasic}}

resource "vkcs_db_backup" "basic" {
    name = "basic"
    dbms_id = vkcs_db_cluster.basic.id
    description = "basic description"
}
`
