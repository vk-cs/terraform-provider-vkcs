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

func TestAccBackupPlan_backup_target(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// creating plan
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanBackupTarget1, map[string]string{
					"testAccBackupPlanBackupTargetBase": acctest.AccTestRenderConfig(testAccBackupPlanBackupTargetBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.#", "1"),
					resource.TestCheckResourceAttrSet("vkcs_backup_plan.basic", "backup_targets.0.instance_id"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.*", "vkcs_blockstorage_volume.bootable", "id"),
				),
			},
			// adding disk
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanBackupTarget2, map[string]string{
					"testAccBackupPlanBackupTargetBase": acctest.AccTestRenderConfig(testAccBackupPlanBackupTargetBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.#", "1"),
					resource.TestCheckResourceAttrSet("vkcs_backup_plan.basic", "backup_targets.0.instance_id"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.*", "vkcs_blockstorage_volume.bootable", "id"),
					resource.TestCheckTypeSetElemAttrPair("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.*", "vkcs_blockstorage_volume.data", "id"),
				),
			},
			// deleting disk
			{
				Config: acctest.AccTestRenderConfig(testAccBackupPlanBackupTarget1, map[string]string{
					"testAccBackupPlanBackupTargetBase": acctest.AccTestRenderConfig(testAccBackupPlanBackupTargetBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.#", "1"),
					resource.TestCheckResourceAttrSet("vkcs_backup_plan.basic", "backup_targets.0.instance_id"),
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("vkcs_backup_plan.basic", "backup_targets.0.volume_ids.*", "vkcs_blockstorage_volume.bootable", "id"),
				),
			},
		},
	})
}

func TestAccBackupPlan_order_independent(t *testing.T) {
	configInstanceIDs := acctest.AccTestRenderConfig(
		testAccBackupPlanInstanceIDsOrder,
		map[string]string{
			"TestAccBackupPlanOrderIndependentBase": acctest.AccTestRenderConfig(testAccBackupPlanOrderIndependentBase),
		},
	)

	configBackupTargets := acctest.AccTestRenderConfig(
		testAccBackupPlanBackupTargetsOrder,
		map[string]string{
			"TestAccBackupPlanOrderIndependentBase": acctest.AccTestRenderConfig(testAccBackupPlanOrderIndependentBase),
		},
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configInstanceIDs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "instance_ids.#", "2"),
				),
			},
			{Config: configInstanceIDs, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configInstanceIDs, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configInstanceIDs, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configInstanceIDs, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configInstanceIDs, PlanOnly: true, ExpectNonEmptyPlan: false},
			{
				Config: configBackupTargets,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_backup_plan.basic", "backup_targets.#", "2"),
				),
			},
			{Config: configBackupTargets, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configBackupTargets, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configBackupTargets, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configBackupTargets, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: configBackupTargets, PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

const testAccBackupPlanBase = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "base_instance" {
  depends_on          = ["vkcs_networking_router_interface.base"]
  name                = "instance_1"
  availability_zone   = "{{.AvailabilityZone}}"
  security_group_ids  = [data.vkcs_networking_secgroup.default_secgroup.id]
  metadata            = {
    foo = "bar"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccBackupPlanRetentionFull = `
{{ .TestAccBackupPlanBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu", "We"]
    time = "11:12+03"
  }
  full_retention     = {
    max_full_backup = 25
  }
  incremental_backup = false
  instance_ids       = [vkcs_compute_instance.base_instance.id]
}
`

const testAccBackupPlanRetentionFullUpdate = `
{{ .TestAccBackupPlanBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu", "We", "Fr"]
    time = "11:15+03"
  }
  full_retention     = {
    max_full_backup = 30
  }
  incremental_backup = false
  instance_ids       = [vkcs_compute_instance.base_instance.id]
}
`

const testAccBackupPlanRetentionGFS = `
{{ .TestAccBackupPlanBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu"]
    time = "16:20"
  }
  gfs_retention      = {
    gfs_weekly  = 10
    gfs_monthly = 2
    gfs_yearly  = 1
  }
  incremental_backup = true
  instance_ids       = [vkcs_compute_instance.base_instance.id]
}
`

const testAccBackupPlanBackupTargetBase = `
{{.BaseNetwork}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}
{{.BaseImage}}

resource "vkcs_blockstorage_volume" "bootable" {
  name              = "bootable-tf-example"
  description       = "test volume"
  metadata          = {
    foo = "bar"
  }
  size              = 10
  availability_zone = "{{.AvailabilityZone}}"
  volume_type       = "ceph-ssd"
  image_id          = data.vkcs_images_image.base.id
}

resource "vkcs_blockstorage_volume" "data" {
  name              = "data-tf-example"
  description       = "test volume"
  metadata          = {
    foo = "bar"
  }
  size              = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type       = "ceph-ssd"
}

resource "vkcs_compute_instance" "base_instance" {
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  flavor_id         = data.vkcs_compute_flavor.base.id
  metadata          = {
    foo = "bar"
  }

  block_device {
    boot_index            = 0
    source_type           = "volume"
    uuid                  = vkcs_blockstorage_volume.bootable.id
    destination_type      = "volume"
    delete_on_termination = true
  }

  block_device {
    boot_index            = -1
    source_type           = "volume"
    uuid                  = vkcs_blockstorage_volume.data.id
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.base.id
  }

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  depends_on = ["vkcs_networking_router_interface.base"]
}
`

const testAccBackupPlanBackupTarget1 = `
{{ .testAccBackupPlanBackupTargetBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu"]
    time = "16:20"
  }
  full_retention     = {
    max_full_backup = 25
  }
  incremental_backup = true
  backup_targets     = [
    {
      instance_id = vkcs_compute_instance.base_instance.id
      volume_ids  = [
        vkcs_blockstorage_volume.bootable.id,
      ]
    }
  ]
}
`

const testAccBackupPlanBackupTarget2 = `
{{ .testAccBackupPlanBackupTargetBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu"]
    time = "16:20"
  }
  full_retention     = {
    max_full_backup = 25
  }
  incremental_backup = true
  backup_targets     = [
    {
      instance_id = vkcs_compute_instance.base_instance.id
      volume_ids  = [
        vkcs_blockstorage_volume.bootable.id,
        vkcs_blockstorage_volume.data.id,
      ]
    }
  ]
}
`

const testAccBackupPlanOrderIndependentBase = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseSecurityGroup}}

data "vkcs_compute_flavor" "small" {
		name = "Basic-1-2-20"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on          = ["vkcs_networking_router_interface.base"]
  name                = "instance_1"
  availability_zone   = "{{.AvailabilityZone}}"
  security_group_ids  = [data.vkcs_networking_secgroup.default_secgroup.id]
  metadata            = {
    foo = "bar"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.small.id
}

resource "vkcs_compute_instance" "instance_2" {
  depends_on          = ["vkcs_networking_router_interface.base"]
  name                = "instance_2"
  availability_zone   = "{{.AvailabilityZone}}"
  security_group_ids  = [data.vkcs_networking_secgroup.default_secgroup.id]
  metadata            = {
    foo = "bar"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.small.id
}
`

const testAccBackupPlanInstanceIDsOrder = `
{{ .TestAccBackupPlanOrderIndependentBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu", "We", "Fr"]
    time = "11:15+03"
  }
  full_retention     = {
    max_full_backup = 30
  }
  incremental_backup = false
  instance_ids       = [
    vkcs_compute_instance.instance_1.id,
    vkcs_compute_instance.instance_2.id,
  ]
}
`

const testAccBackupPlanBackupTargetsOrder = `
{{ .TestAccBackupPlanOrderIndependentBase }}

resource "vkcs_backup_plan" "basic" {
  name               = "tfacc-backup-plan"
  provider_name      = "cloud_servers"
  schedule           = {
    date = ["Tu", "We", "Fr"]
    time = "11:15+03"
  }
  full_retention     = {
    max_full_backup = 30
  }
  incremental_backup = false
  backup_targets     = [
    {
      instance_id = vkcs_compute_instance.instance_1.id
    },
    {
      instance_id = vkcs_compute_instance.instance_2.id
    }
  ]
}
`
