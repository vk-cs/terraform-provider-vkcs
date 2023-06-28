package backup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBackupPlan_basic_retention_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanRetentionFull, map[string]string{
					"TestAccBackupPlanBase": acctest.AccTestRenderConfig(testAccBackupPlanBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "name", "tfacc-backup-plan"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "incremental_backup", "false"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "full_retention.max_full_backup", "25"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.#", "2"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.0", "Tu"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.1", "We"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.time", "11:12+03"),
				),
			},
			{
				ResourceName:            "vkcs_backup_plan.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"schedule.time"},
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanRetentionFullUpdate, map[string]string{
					"TestAccBackupPlanBase": acctest.AccTestRenderConfig(testAccBackupPlanBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "name", "tfacc-backup-plan"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "incremental_backup", "false"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "full_retention.max_full_backup", "30"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.#", "3"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.0", "Tu"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.1", "We"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.2", "Fr"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.time", "11:15+03"),
				),
			},
		},
	})
}

func TestAccBackupPlan_basic_retention_gfs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanRetentionGFS, map[string]string{
					"TestAccBackupPlanBase": acctest.AccTestRenderConfig(testAccBackupPlanBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "name", "tfacc-backup-plan"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "incremental_backup", "true"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "gfs_retention.gfs_weekly", "10"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "gfs_retention.gfs_monthly", "2"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "gfs_retention.gfs_yearly", "1"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.#", "1"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.date.0", "Tu"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "schedule.time", "16:20"),
				),
			},
		},
	})
}

const testAccBackupPlanBase = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "base_instance" {
	depends_on = ["vkcs_networking_router_interface.base"]
	name = "instance_1"
	availability_zone = "{{.AvailabilityZone}}"
	security_groups = ["default"]
	metadata = {
	  foo = "bar"
	}
	network {
	  uuid = vkcs_networking_network.base.id
	}
	image_id = data.vkcs_images_image.base.id
	flavor_id = data.vkcs_compute_flavor.base.id
  }
`

const testAccBackupPlanRetentionFull = `
{{ .TestAccBackupPlanBase }}
resource "vkcs_backup_plan" "basic" {
	name        = "tfacc-backup-plan"
	provider_name = "cloud_servers"
	schedule = {
	  date = ["Tu", "We"]
	  time = "11:12+03"
	}
	full_retention = {
	  max_full_backup = 25
	}
	incremental_backup = false
	instance_ids       = [vkcs_compute_instance.base_instance.id]
  } 
`

const testAccBackupPlanRetentionFullUpdate = `
{{ .TestAccBackupPlanBase }}
resource "vkcs_backup_plan" "basic" {
	name        = "tfacc-backup-plan"
	provider_name = "cloud_servers"
	schedule = {
	  date = ["Tu", "We", "Fr"]
	  time = "11:15+03"
	}
	full_retention = {
	  max_full_backup = 30
	}
	incremental_backup = false
	instance_ids       = [vkcs_compute_instance.base_instance.id]
  } 
`

const testAccBackupPlanRetentionGFS = `
{{ .TestAccBackupPlanBase }}
resource "vkcs_backup_plan" "basic" {
	name        = "tfacc-backup-plan"
	provider_name = "cloud_servers"
	schedule = {
	  date = ["Tu"]
	  time = "16:20"
	}
	gfs_retention = {
		gfs_weekly  = 10
		gfs_monthly = 2
		gfs_yearly = 1
	}
	incremental_backup = true
	instance_ids       = [vkcs_compute_instance.base_instance.id]
  } 
`
