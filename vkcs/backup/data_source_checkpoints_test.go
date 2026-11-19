package backup_test

import (
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBackupCheckpoints_dataSource(t *testing.T) {
	planName := "tfacc-checkpoints-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	config := acctest.AccTestRenderConfig(testAccBackupCheckpointsDataSource, map[string]string{
		"TestAccBackupPlanBase": acctest.AccTestRenderConfig(testAccBackupPlanBase),
		"PlanName":              planName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_backup_checkpoints.plan", "id", "vkcs_backup_plan.checkpoints", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_backup_checkpoints.plan", "plan_id", "vkcs_backup_plan.checkpoints", "id"),
				),
			},
		},
	})
}

const testAccBackupCheckpointsDataSource = `
{{ .TestAccBackupPlanBase }}

resource "vkcs_backup_plan" "checkpoints" {
  name               = "{{.PlanName}}"
  provider_name      = "cloud_servers"
  schedule = {
    every_hours = 24
  }
  full_retention = {
    max_full_backup = 30
  }
  incremental_backup = false
  instance_ids       = [vkcs_compute_instance.base_instance.id]
}

data "vkcs_backup_checkpoints" "plan" {
  plan_id = vkcs_backup_plan.checkpoints.id
}
`
