package compute_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeInstance_checkpointRestore(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("set TF_ACC=1 to run acceptance tests")
	}

	checkpointID := os.Getenv("VKCS_TEST_BACKUP_CHECKPOINT_ID")
	sourceResourceID := os.Getenv("VKCS_TEST_BACKUP_SOURCE_RESOURCE_ID")
	networkID := os.Getenv("VKCS_TEST_BACKUP_NETWORK_ID")
	keyPair := os.Getenv("VKCS_TEST_BACKUP_KEY_PAIR")
	if checkpointID == "" || sourceResourceID == "" || networkID == "" || keyPair == "" {
		t.Skip("set VKCS_TEST_BACKUP_CHECKPOINT_ID, VKCS_TEST_BACKUP_SOURCE_RESOURCE_ID, VKCS_TEST_BACKUP_NETWORK_ID, and VKCS_TEST_BACKUP_KEY_PAIR")
	}

	name := "tfacc-checkpoint-restore-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	config := fmt.Sprintf(`
resource "vkcs_compute_instance" "restored" {
  name               = %q
  checkpoint_id      = %q
  source_resource_id = %q
  key_pair           = %q

  network {
    uuid = %q
  }
}

data "vkcs_backup_restore" "status" {
  restore_id = vkcs_compute_instance.restored.restore_id
}
`, name, checkpointID, sourceResourceID, keyPair, networkID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vkcs_compute_instance.restored", "id"),
					resource.TestCheckResourceAttrSet("vkcs_compute_instance.restored", "restore_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_backup_restore.status", "id"),
					testAccCheckBackupRestoreStatusHasEntries,
				),
			},
		},
	})
}

func testAccCheckBackupRestoreStatusHasEntries(state *terraform.State) error {
	dataSource, ok := state.RootModule().Resources["data.vkcs_backup_restore.status"]
	if !ok || dataSource.Primary == nil {
		return fmt.Errorf("backup restore status data source was not found in state")
	}

	count, err := strconv.Atoi(dataSource.Primary.Attributes["resources_status.%"])
	if err != nil {
		return fmt.Errorf("parse backup restore resource status count: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("backup restore returned no resource statuses")
	}

	return nil
}
