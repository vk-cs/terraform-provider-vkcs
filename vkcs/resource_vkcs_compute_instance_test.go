package vkcs

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestAccComputeInstance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBasic, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "availability_zone", osAvailabilityZone),
				),
			},
		},
	})
}

func TestAccComputeInstance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateActive, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateShutoff, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateActive, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateShutoff, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateActive, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateShutoff, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateActive, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateShelve, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shelved_offloaded"),
					testAccCheckComputeInstanceState(&instance, "shelved_offloaded"),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceStateActive, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBootFromVolumeImage, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBootFromVolumeVolume, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBootFromVolumeForceNew1, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceBootFromVolumeForceNew2, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBlockDeviceNewVolume, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBlockDeviceNewVolumeTypeAndBus, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_blockDeviceExistingVolume(t *testing.T) {
	var instance servers.Server
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceBlockDeviceExistingVolume, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckBlockStorageVolumeExists(
						"vkcs_blockstorage_volume.volume_1", &volume),
				),
			},
		},
	})
}

// TODO: verify the personality really exists on the instance.
func TestAccComputeInstance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstancePersonality, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceAccessIPv4, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceChangeFixedIP1, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceChangeFixedIP2, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceStopBeforeDestroy, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceMetadataRemove1, testAccValues),
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
				Config: testAccRenderConfig(testAccComputeInstanceMetadataRemove2, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceForceDelete, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceTimeout, testAccValues),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceNetworkModeNone, testAccValues),
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
	var network networks.Network
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceNetworkNameToID, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.1.uuid", &network.ID),
				),
			},
		},
	})
}

func TestAccComputeInstance_crazyNICs(t *testing.T) {
	var instance servers.Server
	var network1 networks.Network
	var network2 networks.Network
	var port1 ports.Port
	var port2 ports.Port
	var port3 ports.Port
	var port4 ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceCrazyNICs, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingNetworkExists(
						"vkcs_networking_network.network_1", &network1),
					testAccCheckNetworkingNetworkExists(
						"vkcs_networking_network.network_2", &network2),
					testAccCheckNetworkingPortExists(
						"vkcs_networking_port.port_1", &port1),
					testAccCheckNetworkingPortExists(
						"vkcs_networking_port.port_2", &port2),
					testAccCheckNetworkingPortExists(
						"vkcs_networking_port.port_3", &port3),
					testAccCheckNetworkingPortExists(
						"vkcs_networking_port.port_4", &port4),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.1.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.2.uuid", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.3.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.4.uuid", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.5.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.6.uuid", &network2.ID),
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
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.5.port", &port1.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.6.port", &port2.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.7.port", &port3.ID),
					resource.TestCheckResourceAttrPtr(
						"vkcs_compute_instance.instance_1", "network.8.port", &port4.ID),
				),
			},
		},
	})
}

func TestAccComputeInstance_tags(t *testing.T) {
	var instance servers.Server

	resourceName := "vkcs_compute_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceTagsCreate, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3"}),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceTagsAdd, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3", "tag4"}),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceTagsDelete, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag2", "tag3"}),
				),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceTagsClear, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, nil),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_instance" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" && server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
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
	instance *servers.Server, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range instance.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

func testAccCheckComputeInstanceNoMetadataKey(
	instance *servers.Server, k string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return nil
		}

		for key := range instance.Metadata {
			if k == key {
				return fmt.Errorf("Metadata found: %s", k)
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
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

resource "vkcs_compute_instance" "instance_1" {
	depends_on = ["vkcs_networking_subnet.base"]
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

const testAccComputeInstanceBootFromVolumeImage = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_blockstorage_volume" "vol_1" {
  name = "vol_1"
  size = 5
  image_id = data.vkcs_images_image.base.id
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
  block_device {
    uuid = "${vkcs_blockstorage_volume.vol_1.id}"
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]  
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
  block_device {
    uuid = data.vkcs_images_image.base.id
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    uuid = "${vkcs_blockstorage_volume.volume_1.id}"
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

const testAccComputeInstancePersonality = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.subnet_1"]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    uuid = "${vkcs_networking_network.network_1.id}"
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]

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

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base", "vkcs_networking_subnet.subnet_1"]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    name = "${vkcs_networking_network.network_1.name}"
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id

}
`

const testAccComputeInstanceCrazyNICs = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_networking_network" "network_2" {
  name = "network_2"
}

resource "vkcs_networking_subnet" "subnet_2" {
  name = "subnet_2"
  network_id = "${vkcs_networking_network.network_2.id}"
  cidr = "192.168.2.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.1.103"
  }
}

resource "vkcs_networking_port" "port_2" {
  name = "port_2"
  network_id = "${vkcs_networking_network.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_2.id}"
    ip_address = "192.168.2.103"
  }
}

resource "vkcs_networking_port" "port_3" {
  name = "port_3"
  network_id = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.1.104"
  }
}

resource "vkcs_networking_port" "port_4" {
  name = "port_4"
  network_id = "${vkcs_networking_network.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_2.id}"
    ip_address = "192.168.2.104"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = [
	"vkcs_networking_subnet.base",
    "vkcs_networking_subnet.subnet_1",
    "vkcs_networking_subnet.subnet_2",
    "vkcs_networking_port.port_1",
    "vkcs_networking_port.port_2",
  ]

  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]

  network {
    uuid = vkcs_networking_network.base.id
  }

  network {
    uuid = "${vkcs_networking_network.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
  }

  network {
    uuid = "${vkcs_networking_network.network_2.id}"
    fixed_ip_v4 = "192.168.2.100"
  }

  network {
    uuid = "${vkcs_networking_network.network_1.id}"
    fixed_ip_v4 = "192.168.1.101"
  }

  network {
    uuid = "${vkcs_networking_network.network_2.id}"
    fixed_ip_v4 = "192.168.2.101"
  }

  network {
    port = "${vkcs_networking_port.port_1.id}"
  }

  network {
    port = "${vkcs_networking_port.port_2.id}"
  }

  network {
    port = "${vkcs_networking_port.port_3.id}"
  }

  network {
    port = "${vkcs_networking_port.port_4.id}"
  }

  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceStateActive = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
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

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`
