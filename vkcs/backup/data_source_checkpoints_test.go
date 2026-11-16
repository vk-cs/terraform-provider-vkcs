package backup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBackupCheckpoints_dataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("set TF_ACC=1 to run acceptance tests")
	}
	planID := os.Getenv("VKCS_TEST_BACKUP_PLAN_ID")
	if planID == "" {
		t.Skip("set VKCS_TEST_BACKUP_PLAN_ID to an existing backup plan")
	}

	config := fmt.Sprintf(`
data "vkcs_backup_checkpoints" "plan" {
  plan_id = %q
}
`, planID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_backup_checkpoints.plan", "id", planID),
					resource.TestCheckResourceAttr("data.vkcs_backup_checkpoints.plan", "plan_id", planID),
				),
			},
		},
	})
}
