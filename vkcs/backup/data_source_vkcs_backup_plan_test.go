package backup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBackupPlanDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config: acctest.AccTestRenderConfig(testAccBackupPlanDataSourceBasic, map[string]string{
					"TestAccBackupPlanDataSourceBase": acctest.AccTestRenderConfig(testAccBackupPlanRetentionFull, map[string]string{
						"TestAccBackupPlanBase": acctest.AccTestRenderConfig(testAccBackupPlanBase),
					})}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "name", "tfacc-backup-plan"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "incremental_backup", "false"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "full_retention.max_full_backup", "25"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "schedule.date.#", "2"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "schedule.date.0", "Tu"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "schedule.date.1", "We"),
					resource.TestCheckResourceAttr("data.vkcs_backup_plan.plan", "schedule.time", "8:12"),
				),
			},
		},
	})
}

const testAccBackupPlanDataSourceBasic = `
{{ .TestAccBackupPlanDataSourceBase}}

data "vkcs_backup_plan" "plan" {
	name = vkcs_backup_plan.basic.name
}
`
