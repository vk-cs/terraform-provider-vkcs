package compute_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"
	"time"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	backupcheckpoints "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	backupproviders "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/providers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccComputeInstance_checkpointRestore(t *testing.T) {
	suffix := acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	sourceName := "tfacc-checkpoint-source-" + suffix
	planName := "tfacc-checkpoint-plan-" + suffix
	restoredName := "tfacc-checkpoint-restore-" + suffix
	keyPairName := "tfacc-restore-key-" + suffix
	fixtureConfig := acctest.AccTestRenderConfig(testAccComputeInstanceCheckpointFixture, map[string]string{
		"SourceName": sourceName,
		"PlanName":   planName,
	})
	config := acctest.AccTestRenderConfig(testAccComputeInstanceCheckpointRestore, map[string]string{
		"TestAccComputeInstanceCheckpointFixture": fixtureConfig,
		"RestoredName": restoredName,
		"KeyPairName":  keyPairName,
		"PublicKey":    testAccGenerateEd25519PublicKey(t),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fixtureConfig,
				Check:  testAccCreateAndWaitForNovaBackupCheckpoint(t, 10*time.Minute),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_compute_instance.restored", "name", restoredName),
					resource.TestCheckResourceAttrSet("vkcs_compute_instance.restored", "id"),
					resource.TestCheckResourceAttrSet("vkcs_compute_instance.restored", "restore_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_backup_restore.status", "id"),
					testAccCheckRestoredInstanceIsDistinctFromSource("vkcs_compute_instance.checkpoint_source"),
					testAccCheckBackupRestoreStatusHasEntries,
				),
			},
		},
	})
}

const testAccComputeInstanceCheckpointRestore = `
{{.TestAccComputeInstanceCheckpointFixture}}

resource "vkcs_compute_keypair" "restore_fixture" {
  name       = "{{.KeyPairName}}"
  public_key = "{{.PublicKey}}"
}

data "vkcs_backup_checkpoints" "source" {
  plan_id = vkcs_backup_plan.checkpoint_source.id
}

locals {
  nova_checkpoint_ids = flatten([
    for protected_resource in data.vkcs_backup_checkpoints.source.resources : [
      for backup in protected_resource.backups : backup.checkpoint_id
      if backup.status == "available"
    ] if protected_resource.resource_type == "OS::Nova::Server"
  ])
}

resource "vkcs_compute_instance" "restored" {
  name               = "{{.RestoredName}}"
  checkpoint_id      = one(local.nova_checkpoint_ids)
  source_resource_id = vkcs_compute_instance.checkpoint_source.id
  key_pair           = vkcs_compute_keypair.restore_fixture.name

  network {
    uuid = vkcs_networking_network.base.id
  }
}

data "vkcs_backup_restore" "status" {
  restore_id = vkcs_compute_instance.restored.restore_id
}
`

const testAccComputeInstanceCheckpointFixture = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "checkpoint_source" {
  depends_on         = [vkcs_networking_router_interface.base]
  name               = "{{.SourceName}}"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  image_id           = data.vkcs_images_image.base.id
  flavor_id          = data.vkcs_compute_flavor.base.id

  network {
    uuid = vkcs_networking_network.base.id
  }
}

resource "vkcs_backup_plan" "checkpoint_source" {
  name               = "{{.PlanName}}"
  provider_name      = "cloud_servers"
  schedule = {
    every_hours = 24
  }
  full_retention = {
    max_full_backup = 2
  }
  incremental_backup = false
  instance_ids       = [vkcs_compute_instance.checkpoint_source.id]
}
`

func testAccGenerateEd25519PublicKey(t *testing.T) string {
	t.Helper()

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate acceptance SSH key: %v", err)
	}
	encodedKey := append(testAccSSHString([]byte("ssh-ed25519")), testAccSSHString(publicKey)...)
	return "ssh-ed25519 " + base64.StdEncoding.EncodeToString(encodedKey) + " terraform-acceptance"
}

func testAccSSHString(value []byte) []byte {
	encoded := make([]byte, 4+len(value))
	binary.BigEndian.PutUint32(encoded, uint32(len(value)))
	copy(encoded[4:], value)
	return encoded
}

func testAccCreateAndWaitForNovaBackupCheckpoint(t *testing.T, timeout time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		plan, ok := state.RootModule().Resources["vkcs_backup_plan.checkpoint_source"]
		if !ok || plan.Primary == nil || plan.Primary.ID == "" {
			return fmt.Errorf("checkpoint source backup plan was not found in state")
		}
		source, ok := state.RootModule().Resources["vkcs_compute_instance.checkpoint_source"]
		if !ok || source.Primary == nil || source.Primary.ID == "" {
			return fmt.Errorf("checkpoint source VM was not found in state")
		}

		config, err := (&clients.ConfigOpts{}).LoadAndValidate()
		if err != nil {
			return fmt.Errorf("load acceptance provider config: %w", err)
		}
		backupClient, err := config.BackupV1Client(config.GetRegion(), config.GetProjectID())
		if err != nil {
			return fmt.Errorf("create acceptance backup client: %w", err)
		}
		allProviders, err := backupproviders.List(backupClient).Extract()
		if err != nil {
			return fmt.Errorf("list acceptance backup providers: %w", err)
		}
		providerID := ""
		for _, backupProvider := range allProviders {
			if backupProvider.Name == "OS::Nova" {
				providerID = backupProvider.ID
				break
			}
		}
		if providerID == "" {
			return fmt.Errorf("Nova backup provider was not found")
		}

		created, err := backupcheckpoints.CreateForPlan(backupClient, providerID, &backupcheckpoints.PlanCheckpointCreate{
			PlanCheckpoints: &backupcheckpoints.PlanCheckpointCreateOpts{PlanID: plan.Primary.ID},
		}).Extract()
		if err != nil {
			return fmt.Errorf("start on-demand checkpoint for plan %q: %w", plan.Primary.ID, err)
		}
		if len(created.Checkpoints) == 0 {
			return fmt.Errorf("on-demand checkpoint API returned no checkpoints for plan %q", plan.Primary.ID)
		}
		checkpointIDs := make(map[string]struct{}, len(created.Checkpoints))
		for _, checkpoint := range created.Checkpoints {
			if checkpoint.ID != "" {
				checkpointIDs[checkpoint.ID] = struct{}{}
			}
		}
		if len(checkpointIDs) == 0 {
			return fmt.Errorf("on-demand checkpoint API returned checkpoints without IDs for plan %q", plan.Primary.ID)
		}
		t.Cleanup(func() {
			for checkpointID := range checkpointIDs {
				if err := backupcheckpoints.Delete(backupClient, providerID, checkpointID).ExtractErr(); err != nil && !errutil.IsNotFound(err) {
					t.Errorf("delete acceptance checkpoint %q: %v", checkpointID, err)
				}
			}
		})

		deadline := time.Now().Add(timeout)
		for {
			response, err := backupcheckpoints.ListGroupedByPlanID(backupClient, backupcheckpoints.ListOpts{PlanID: plan.Primary.ID}).Extract()
			if err != nil {
				return fmt.Errorf("read checkpoint status for plan %q: %w", plan.Primary.ID, err)
			}
			for _, protectedResource := range response.Resources {
				if protectedResource.ResourceType != "OS::Nova::Server" {
					continue
				}
				for _, backup := range protectedResource.Backups {
					if _, requested := checkpointIDs[backup.CheckpointID]; !requested {
						continue
					}
					if backup.Status == "available" {
						return nil
					}
					if backup.Status == "error" {
						return fmt.Errorf("on-demand checkpoint %q failed for plan %q", backup.CheckpointID, plan.Primary.ID)
					}
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("timed out waiting for on-demand Nova checkpoint for plan %q", plan.Primary.ID)
			}
			time.Sleep(15 * time.Second)
		}
	}
}

func testAccCheckRestoredInstanceIsDistinctFromSource(sourceResource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instance, ok := state.RootModule().Resources["vkcs_compute_instance.restored"]
		if !ok || instance.Primary == nil {
			return fmt.Errorf("restored compute instance was not found in state")
		}
		source, ok := state.RootModule().Resources[sourceResource]
		if !ok || source.Primary == nil {
			return fmt.Errorf("checkpoint source VM was not found in state")
		}
		if instance.Primary.ID == source.Primary.ID {
			return fmt.Errorf("restored VM ID %q matches the source VM ID", instance.Primary.ID)
		}
		return nil
	}
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
