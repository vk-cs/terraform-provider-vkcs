package compute_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iservers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servers"
)

func TestAccComputeInstance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "availability_zone", acctest.OsAvailabilityZone),
				),
			},
		},
	})
}

func TestAccComputeInstance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateActive),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateShutoff),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateActive),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeInstance_initialStateShutoff(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateShutoff),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateActive),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateShutoff),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
		},
	})
}

func TestAccComputeInstance_initialShelve(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateActive),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateShelve),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shelved_offloaded"),
					testAccCheckComputeInstanceState(&instance, "shelved_offloaded"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStateActive),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeInstance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBootFromVolumeImage),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBootFromVolumeVolume),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_bootFromVolumeForceNew(t *testing.T) {
	var instance1 servers.Server
	var instance2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBootFromVolumeForceNew1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBootFromVolumeForceNew2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance2),
					testAccCheckComputeInstanceInstanceIDsDoNotMatch(&instance1, &instance2),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceNewVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBlockDeviceNewVolume),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceNewVolumeTypeAndBus(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBlockDeviceNewVolumeTypeAndBus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceExistingVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBlockDeviceExistingVolume),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceBootIndexDefault(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBlockDeviceBootIndexDefault),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceBootIndicesDefaults(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceBlockDevicesBootIndicesDefaults),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr("vkcs_compute_instance.instance_1", "block_device.0.boot_index", "0"),
					resource.TestCheckResourceAttr("vkcs_compute_instance.instance_1", "block_device.1.boot_index", "-1"),
					resource.TestCheckResourceAttr("vkcs_compute_instance.instance_1", "block_device.2.boot_index", "-1"),
				),
			},
		},
	})
}

// TODO: verify the personality really exists on the instance.
func TestAccComputeInstance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstancePersonality),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_accessIPv4(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceAccessIPv4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "access_ip_v4", "192.168.1.100"),
				),
			},
		},
	})
}

func TestAccComputeInstance_changeFixedIP(t *testing.T) {
	var instance1 servers.Server
	var instance2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceChangeFixedIP1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceChangeFixedIP2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance2),
					testAccCheckComputeInstanceInstanceIDsDoNotMatch(&instance1, &instance2),
				),
			},
		},
	})
}

func TestAccComputeInstance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceStopBeforeDestroy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_metadataRemove(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceMetadataRemove1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.abc", "def"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceMetadataRemove2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeInstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeInstance_forceDelete(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceForceDelete),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceTimeout),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_networkModeNone(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceNetworkModeNone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceNetworkDoesNotExist("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_networkNameToID(t *testing.T) {
	var instance servers.Server
	uniqueFields := acctest.GenerateUniqueTestFields(t.Name())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceNetworkNameToID, uniqueFields),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.1.uuid", "vkcs_networking_network.network_acc_test", "id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_crazyNICs(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceCrazyNICs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.1.uuid", "vkcs_networking_network.network_1", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.2.uuid", "vkcs_networking_network.network_2", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.3.uuid", "vkcs_networking_network.network_1", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.4.uuid", "vkcs_networking_network.network_2", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.5.uuid", "vkcs_networking_network.network_1", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.6.uuid", "vkcs_networking_network.network_2", "id"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.1.name", "network_1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.2.name", "network_2"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.3.name", "network_1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.4.name", "network_2"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.5.name", "network_1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.6.name", "network_2"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.7.name", "network_1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.8.name", "network_2"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.1.fixed_ip_v4", "192.168.1.100"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.2.fixed_ip_v4", "192.168.2.100"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.3.fixed_ip_v4", "192.168.1.101"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "network.4.fixed_ip_v4", "192.168.2.101"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.5.port", "vkcs_networking_port.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.6.port", "vkcs_networking_port.port_2", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.7.port", "vkcs_networking_port.port_3", "id"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "network.8.port", "vkcs_networking_port.port_4", "id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_tags(t *testing.T) {
	var instance servers.Server

	resourceName := "vkcs_compute_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceTagsCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceTagsAdd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3", "tag4"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceTagsDelete),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag2", "tag3"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceTagsClear),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, nil),
				),
			},
		},
	})
}

func TestResourceInstanceBlockDevicesV2(t *testing.T) {
	for _, c := range []struct {
		name         string
		resourceData *schema.ResourceData
		blockDevices []map[string]interface{}
		expected     []bootfromvolume.BlockDevice
		expectedErr  error
	}{
		{
			"one block device",
			nil,
			[]map[string]interface{}{
				{
					"uuid":                  "7a619268-4432-456a-a67c-0f5c58933c2c",
					"volume_size":           10,
					"boot_index":            -1,
					"delete_on_termination": true,
					"guest_format":          "ext2",
					"volume_type":           "ceph-ssd",
					"device_type":           "disk",
					"disk_bus":              "virtio",
					"source_type":           "image",
					"destination_type":      "local",
				},
			},
			[]bootfromvolume.BlockDevice{
				{
					UUID:                "7a619268-4432-456a-a67c-0f5c58933c2c",
					VolumeSize:          10,
					BootIndex:           0,
					DeleteOnTermination: true,
					GuestFormat:         "ext2",
					VolumeType:          "ceph-ssd",
					DeviceType:          "disk",
					DiskBus:             "virtio",
					SourceType:          bootfromvolume.SourceImage,
					DestinationType:     bootfromvolume.DestinationLocal,
				},
			},
			nil,
		},

		{
			"multiple block devices",
			nil,
			[]map[string]interface{}{
				{
					"uuid":                  "7a619268-4432-456a-a67c-0f5c58933c2c",
					"volume_size":           10,
					"boot_index":            -1,
					"delete_on_termination": true,
					"guest_format":          "ext2",
					"volume_type":           "ceph-ssd",
					"device_type":           "disk",
					"disk_bus":              "virtio",
					"source_type":           "image",
					"destination_type":      "local",
				},
				{
					"uuid":                  "44338d06-ff25-4cc2-bedd-df3c1a8cab28",
					"volume_size":           10,
					"boot_index":            -1,
					"delete_on_termination": true,
					"guest_format":          "ext2",
					"volume_type":           "ceph-ssd",
					"device_type":           "disk",
					"disk_bus":              "virtio",
					"source_type":           "image",
					"destination_type":      "local",
				},
				{
					"uuid":                  "34d137c5-0d8e-4db2-bae1-1f6d394a9e52",
					"volume_size":           10,
					"boot_index":            0,
					"delete_on_termination": true,
					"guest_format":          "ext2",
					"volume_type":           "ceph-ssd",
					"device_type":           "disk",
					"disk_bus":              "virtio",
					"source_type":           "image",
					"destination_type":      "local",
				},
			},
			[]bootfromvolume.BlockDevice{
				{
					UUID:                "7a619268-4432-456a-a67c-0f5c58933c2c",
					VolumeSize:          10,
					BootIndex:           -1,
					DeleteOnTermination: true,
					GuestFormat:         "ext2",
					VolumeType:          "ceph-ssd",
					DeviceType:          "disk",
					DiskBus:             "virtio",
					SourceType:          bootfromvolume.SourceImage,
					DestinationType:     bootfromvolume.DestinationLocal,
				},
				{
					UUID:                "44338d06-ff25-4cc2-bedd-df3c1a8cab28",
					VolumeSize:          10,
					BootIndex:           -1,
					DeleteOnTermination: true,
					GuestFormat:         "ext2",
					VolumeType:          "ceph-ssd",
					DeviceType:          "disk",
					DiskBus:             "virtio",
					SourceType:          bootfromvolume.SourceImage,
					DestinationType:     bootfromvolume.DestinationLocal,
				},
				{
					UUID:                "34d137c5-0d8e-4db2-bae1-1f6d394a9e52",
					VolumeSize:          10,
					BootIndex:           0,
					DeleteOnTermination: true,
					GuestFormat:         "ext2",
					VolumeType:          "ceph-ssd",
					DeviceType:          "disk",
					DiskBus:             "virtio",
					SourceType:          bootfromvolume.SourceImage,
					DestinationType:     bootfromvolume.DestinationLocal,
				},
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			bds := make([]interface{}, len(c.blockDevices))
			for i, bd := range c.blockDevices {
				bds[i] = bd
			}
			actual, actualErr := compute.ResourceInstanceBlockDevicesV2(c.resourceData, bds)
			if !errors.Is(c.expectedErr, actualErr) {
				t.Fatalf("Unexpected error: wanted %s got %s", c.expectedErr, actualErr)
			}
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestAccComputeInstance_cloudMonitoringBasic(t *testing.T) {
	var instance servers.Server
	resourceName := "vkcs_compute_instance.instance_1"
	monitoringResourceName := "vkcs_cloud_monitoring.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceCloudMonitoring, map[string]string{"Metadata": testMetadataEmpty}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttrSet(monitoringResourceName, "service_user_id"),
					resource.TestCheckResourceAttrSet(monitoringResourceName, "script"),
					resource.TestCheckResourceAttrPair(
						resourceName, "all_metadata.service_user_id", monitoringResourceName, "service_user_id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_cloudMonitoringMetadata(t *testing.T) {
	var instance servers.Server
	resourceName := "vkcs_compute_instance.instance_1"
	monitoringResourceName := "vkcs_cloud_monitoring.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceCloudMonitoring, map[string]string{"Metadata": testMetadata}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceMetadata(&instance, "key1", "value1"),
					testAccCheckComputeInstanceMetadata(&instance, "key2", "value2"),
					resource.TestCheckResourceAttrPair(
						resourceName, "all_metadata.service_user_id", monitoringResourceName, "service_user_id"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceCloudMonitoring, map[string]string{"Metadata": testMetadataUpdate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceMetadata(&instance, "key1", "value1"),
					testAccCheckComputeInstanceMetadata(&instance, "key3", "value3"),
					testAccCheckComputeInstanceNoMetadataKey(&instance, "key2"),
					resource.TestCheckResourceAttrPair(
						resourceName, "all_metadata.service_user_id", monitoringResourceName, "service_user_id"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceCloudMonitoring, map[string]string{"Metadata": testMetadataEmpty}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceNoMetadataKey(&instance, "key1"),
					testAccCheckComputeInstanceNoMetadataKey(&instance, "key3"),
					resource.TestCheckResourceAttrPair(
						resourceName, "all_metadata.service_user_id", monitoringResourceName, "service_user_id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_GetWindowsPassword(t *testing.T) {
	var instance servers.Server
	resourceName := "vkcs_compute_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceGetWindowsPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttrSet(resourceName, "password_data"),
				),
			},
		},
	})
}

func TestAccComputeInstance_AdminPasswordUpdate(t *testing.T) {
	var instance servers.Server
	resourceName := "vkcs_compute_instance.instance_1"
	oldPassword := "S1mfd)dIt8s_mD"
	newPassword := "foDp2x_f(sdoRv"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceAdminPassword, map[string]string{"AdminPass": oldPassword}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "admin_pass", oldPassword),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceAdminPassword, map[string]string{"AdminPass": newPassword}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "admin_pass", newPassword),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	config, err := clients.ConfigureFromEnv(context.Background())
	if err != nil {
		return fmt.Errorf("Error authenticating clients from environment: %s", err)
	}

	computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_instance" {
			continue
		}

		server, err := iservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" && server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(resourceName string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config, err := clients.ConfigureFromEnv(context.Background())
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := iservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeInstanceMetadata(
	instance *servers.Server, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual, ok := instance.Metadata[key]
		if !ok {
			return fmt.Errorf("Metadata not found: %s", key)
		}
		if value != actual {
			return fmt.Errorf("Bad value for key %s, expected: %s, actual: %s", key, value, actual)
		}

		return nil
	}
}

func testAccCheckComputeInstanceNoMetadataKey(
	instance *servers.Server, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, ok := instance.Metadata[key]; ok {
			return fmt.Errorf("Unexpected metadata found: %s", key)
		}

		return nil
	}
}

func testAccCheckComputeInstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {
				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("Unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})
		if err != nil {
			return fmt.Errorf("Unable to list volume attachments")
		}

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("No attached volume found")
	}
}

func testAccCheckComputeInstanceInstanceIDsDoNotMatch(
	instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("Instance was not recreated")
		}

		return nil
	}
}

func testAccCheckComputeInstanceState(
	instance *servers.Server, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.ToLower(instance.Status) != state {
			return fmt.Errorf("Instance state is not match")
		}

		return nil
	}
}

func testAccCheckComputeInstanceTags(name string, tags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		if _, ok := rs.Primary.Attributes["tags.#"]; !ok {
			return fmt.Errorf("resource tags not found: %s.tags", name)
		}

		var rtags []string
		for key, val := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "tags.") {
				continue
			}

			if key == "tags.#" {
				continue
			}

			rtags = append(rtags, val)
		}

		sort.Strings(rtags)
		sort.Strings(tags)
		if !reflect.DeepEqual(rtags, tags) {
			return fmt.Errorf(
				"%s.tags: expected: %#v, got %#v", name, tags, rtags)
		}
		return nil
	}
}

func testAccCheckComputeInstanceNetworkDoesNotExist(n string, _ *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "0" {
			return fmt.Errorf("network should not exists when network mode 'none': %s", n)
		}

		return nil
	}
}

const testAccComputeInstanceBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
	depends_on = ["vkcs_networking_router_interface.base"]
	name = "instance_1"
	availability_zone = "{{.AvailabilityZone}}"
	security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
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

const testAccComputeInstanceBootFromVolumeImage = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBootFromVolumeVolume = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "vol_1" {
  name = "vol_1"
  size = 5
  image_id = data.vkcs_images_image.base.id
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = vkcs_blockstorage_volume.vol_1.id
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBootFromVolumeForceNew1 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]  
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBootFromVolumeForceNew2 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    volume_size = 4
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBlockDeviceNewVolume = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBlockDeviceNewVolumeTypeAndBus = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    destination_type = "local"
    boot_index = 0
	delete_on_termination = true
	device_type = "disk"
	disk_bus = "virtio"
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
	delete_on_termination = true
	device_type = "disk"
	disk_bus = "virtio"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`
const testAccComputeInstanceBlockDeviceExistingVolume = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    uuid = vkcs_blockstorage_volume.volume_1.id
    source_type = "volume"
    destination_type = "volume"
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBlockDeviceBootIndexDefault = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    destination_type = "local"
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceBlockDevicesBootIndicesDefaults = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
	boot_index = 0
    destination_type = "local"
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstancePersonality = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  personality {
    file = "/tmp/foobar.txt"
    content = "happy"
  }
  personality {
    file = "/tmp/barfoo.txt"
    content = "angry"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceAccessIPv4 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = vkcs_networking_network.network_1.id
  cidr = "192.168.1.0/24"
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.subnet_1"]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    uuid = vkcs_networking_network.network_1.id
    fixed_ip_v4 = "192.168.1.100"
    access_network = true
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceChangeFixedIP1 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
    fixed_ip_v4 = "192.168.199.134"
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceChangeFixedIP2 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
    fixed_ip_v4 = "192.168.199.135"
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceStopBeforeDestroy = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  stop_before_destroy = true
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceMetadataRemove1 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  metadata = {
    foo = "bar"
    abc = "def"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceMetadataRemove2 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceForceDelete = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  force_delete = true
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceTimeout = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  timeouts {
    create = "10m"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceNetworkModeNone = `
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
  name = "test-instance-1"
  availability_zone = "{{.AvailabilityZone}}"

  network_mode = "none"
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceNetworkNameToID = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_networking_network" "network_acc_test" {
  name = "network_{{.TestName}}_{{.CurrentTime}}"
}

resource "vkcs_networking_subnet" "subnet_acc_test" {
  name = "subnet_acc_test"
  network_id = vkcs_networking_network.network_acc_test.id
  cidr = "192.168.1.0/24"
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base", "vkcs_networking_subnet.subnet_acc_test"]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    name = vkcs_networking_network.network_acc_test.name
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceCrazyNICs = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = vkcs_networking_network.network_1.id
  cidr = "192.168.1.0/24"
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_networking_network" "network_2" {
  name = "network_2"
}

resource "vkcs_networking_subnet" "subnet_2" {
  name = "subnet_2"
  network_id = vkcs_networking_network.network_2.id
  cidr = "192.168.2.0/24"
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = vkcs_networking_network.network_1.id
  admin_state_up = "true"

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.1.103"
  }
}

resource "vkcs_networking_port" "port_2" {
  name = "port_2"
  network_id = vkcs_networking_network.network_2.id
  admin_state_up = "true"

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_2.id
    ip_address = "192.168.2.103"
  }
}

resource "vkcs_networking_port" "port_3" {
  name = "port_3"
  network_id = vkcs_networking_network.network_1.id
  admin_state_up = "true"

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.1.104"
  }
}

resource "vkcs_networking_port" "port_4" {
  name = "port_4"
  network_id = vkcs_networking_network.network_2.id
  admin_state_up = "true"

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_2.id
    ip_address = "192.168.2.104"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = [
	"vkcs_networking_router_interface.base",
    "vkcs_networking_subnet.subnet_1",
    "vkcs_networking_subnet.subnet_2",
    "vkcs_networking_port.port_1",
    "vkcs_networking_port.port_2",
  ]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    uuid = vkcs_networking_network.network_1.id
    fixed_ip_v4 = "192.168.1.100"
  }

  network {
    uuid = vkcs_networking_network.network_2.id
    fixed_ip_v4 = "192.168.2.100"
  }

  network {
    uuid = vkcs_networking_network.network_1.id
    fixed_ip_v4 = "192.168.1.101"
  }

  network {
    uuid = vkcs_networking_network.network_2.id
    fixed_ip_v4 = "192.168.2.101"
  }

  network {
    port = vkcs_networking_port.port_1.id
  }

  network {
    port = vkcs_networking_port.port_2.id
  }

  network {
    port = vkcs_networking_port.port_3.id
  }

  network {
    port = vkcs_networking_port.port_4.id
  }

  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceStateActive = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = [vkcs_networking_router_interface.base]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  power_state = "active"
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceStateShutoff = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  power_state = "shutoff"
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceStateShelve = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  power_state = "shelved_offloaded"
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceTagsCreate = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  tags = ["tag1", "tag2", "tag3"]
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`
const testAccComputeInstanceTagsAdd = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  tags = ["tag1", "tag2", "tag3", "tag4"]
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceTagsDelete = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  tags = ["tag2", "tag3"]
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceTagsClear = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceCloudMonitoring = `
{{.BaseNetwork}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

data "vkcs_images_image" "base" {
  visibility = "public"
  default    = true
  properties = {
    mcs_os_distro  = "ubuntu"
    mcs_os_version = "24.04"
  }
}

resource "vkcs_cloud_monitoring" "basic" {
  image_id = data.vkcs_images_image.base.id
}

resource "vkcs_compute_instance" "instance_1" {
  name               = "instance-with-monitoring"
  image_id           = data.vkcs_images_image.base.id
  flavor_id          = data.vkcs_compute_flavor.base.id
  availability_zone  = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }

  cloud_monitoring {
    service_user_id = vkcs_cloud_monitoring.basic.service_user_id
    script          = vkcs_cloud_monitoring.basic.script
  }

  {{.Metadata}}

  user_data = <<EOF
#cloud-config
runcmd:
  - echo "Hello, world!" > /etc/motd
  - [ sh, -c, "echo 'Second command executed successfully!' >> /run/testing.txt" ]
write_files:
  - path: /etc/example_config.conf
    content: |
      [example-config]
      key=value
final_message: "The system is up, after $UPTIME seconds"
  EOF

  depends_on = [
    vkcs_networking_router_interface.base
  ]
}
`

const testMetadata = `
metadata = {
  key1 = "value1"
  key2 = "value2"
}
`

const testMetadataUpdate = `
metadata = {
  key1 = "value1"
  key3 = "value3"
}
`

const testMetadataEmpty = `
metadata = {}
`

const testAccComputeInstanceGetWindowsPassword = `
{{.BaseNetwork}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

data "vkcs_images_image" "windows" {
  visibility  = "public"
  default     = true
  most_recent = true
  properties = {
    mcs_os_type = "windows"
    os_version  = "10.0"
  }
}

resource "vkcs_compute_keypair" "windows" {
  name = "windows-key-pair"
}

resource "vkcs_compute_instance" "instance_1" {
  name               = "instance-with-windows"
  flavor_id          = data.vkcs_compute_flavor.base.id
  availability_zone  = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  key_pair 			 = vkcs_compute_keypair.windows.name
  network {
    uuid = vkcs_networking_network.base.id
  }

  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.windows.id
    destination_type      = "volume"
    volume_size           = 50
    volume_type           = "ceph-ssd"
    delete_on_termination = true
  }

  vendor_options {
    get_password_data = true
  }
  
  depends_on = [
    vkcs_networking_router_interface.base
  ]
}
`

const testAccComputeInstanceAdminPassword = `
{{.BaseNetwork}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

data "vkcs_images_image" "base" {
  visibility = "public"
  default    = true
  properties = {
    mcs_os_distro  = "ubuntu"
    mcs_os_version = "24.04"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  name               = "instance-with-admin-password"
  image_id           = data.vkcs_images_image.base.id
  flavor_id          = data.vkcs_compute_flavor.base.id
  availability_zone  = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  admin_pass         = "{{.AdminPass}}"
  network {
    uuid = vkcs_networking_network.base.id
  }

  depends_on = [
    vkcs_networking_router_interface.base
  ]
}
`
