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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "availability_zone", "MS1"),
				),
			},
		},
	})
}

func TestAccComputeInstance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeInstanceStateShutoff(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeInstanceStateActive(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceStateShutoff(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shutoff"),
					testAccCheckComputeInstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeInstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeInstanceStateShutoff(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "active"),
					testAccCheckComputeInstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeInstanceStateShelve(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_compute_instance.instance_1", "power_state", "shelved_offloaded"),
					testAccCheckComputeInstanceState(&instance, "shelved_offloaded"),
				),
			},
			{
				Config: testAccComputeInstanceStateActive(),
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

// func TestAccComputeInstance_secgroupMulti(t *testing.T) {
// 	var instance1 servers.Server
// 	var secgroup1 secgroups.SecurityGroup

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckCompute(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckComputeInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccComputeInstanceSecgroupMulti(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckComputeSecGroupExists(
// 						"vkcs_compute_secgroup.secgroup_1", &secgroup1),
// 					testAccCheckComputeInstanceExists(
// 						"vkcs_compute_instance.instance_1", &instance1),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccComputeInstance_secgroupMultiUpdate(t *testing.T) {
// 	var instance1 servers.Server
// 	var secgroup1, secgroup2 secgroups.SecurityGroup

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckCompute(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckComputeInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccComputeInstanceSecgroupMultiUpdate1(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckComputeSecGroupExists(
// 						"vkcs_compute_secgroup.secgroup_1", &secgroup1),
// 					testAccCheckComputeSecGroupExists(
// 						"vkcs_compute_secgroup.secgroup_2", &secgroup2),
// 					testAccCheckComputeInstanceExists(
// 						"vkcs_compute_instance.instance_1", &instance1),
// 				),
// 			},
// 			{
// 				Config: testAccComputeInstanceSecgroupMultiUpdate2(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckComputeSecGroupExists(
// 						"vkcs_compute_secgroup.secgroup_1", &secgroup1),
// 					testAccCheckComputeSecGroupExists(
// 						"vkcs_compute_secgroup.secgroup_2", &secgroup2),
// 					testAccCheckComputeInstanceExists(
// 						"vkcs_compute_instance.instance_1", &instance1),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccComputeInstance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBootFromVolumeImage(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBootFromVolumeVolume(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBootFromVolumeForceNew1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: testAccComputeInstanceBootFromVolumeForceNew2(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBlockDeviceNewVolume(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBlockDeviceNewVolumeTypeAndBus(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceBlockDeviceExistingVolume(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckBlockStorageVolumeExists(
						"vkcs_blockstorage_volume.volume_1", &volume),
				),
			},
		},
	})
}

//TODO: verify the personality really exists on the instance.
func TestAccComputeInstance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancePersonality(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

// func TestAccComputeInstance_multiEphemeral(t *testing.T) {
// 	var instance servers.Server

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckCompute(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckComputeInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccComputeInstanceMultiEphemeral(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckComputeInstanceExists(
// 						"vkcs_compute_instance.instance_1", &instance),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccComputeInstance_accessIPv4(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceAccessIPv4(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceChangeFixedIP1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						"vkcs_compute_instance.instance_1", &instance1),
				),
			},
			{
				Config: testAccComputeInstanceChangeFixedIP2(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceStopBeforeDestroy(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceMetadataRemove1(),
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
				Config: testAccComputeInstanceMetadataRemove2(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceForceDelete(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTimeout(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

// func TestAccComputeInstance_networkModeAuto(t *testing.T) {
// 	var instance servers.Server
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckCompute(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckComputeInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccComputeInstanceNetworkModeAuto(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
// 					testAccCheckComputeInstanceNetworkExists("vkcs_compute_instance.instance_1", &instance),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccComputeInstance_networkModeNone(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceNetworkModeNone(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceNetworkNameToID(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceCrazyNICs(),
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
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTagsCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeInstanceTagsAdd(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag1", "tag2", "tag3", "tag4"}),
				),
			},
			{
				Config: testAccComputeInstanceTagsDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(resourceName, []string{"tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeInstanceTagsClear(),
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
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
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
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
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

func testAccCheckComputeInstanceNetworkExists(n string, _ *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "1" {
			return fmt.Errorf("network should be exists when network mode 'auto': %s", n)
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

func testAccComputeInstanceBasic() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  availability_zone = "MS1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

// func testAccComputeInstanceSecgroupMulti() string {
// 	return fmt.Sprintf(`
// resource "vkcs_compute_secgroup" "secgroup_1" {
//   name = "secgroup_1"
//   description = "a security group"
//   rule {
//     from_port = 22
//     to_port = 22
//     ip_protocol = "tcp"
//     cidr = "0.0.0.0/0"
//   }
// }

// resource "vkcs_compute_instance" "instance_1" {
//   name = "instance_1"
//   security_groups = ["default", "${vkcs_compute_secgroup.secgroup_1.name}"]
//   network {
//     uuid = "%s"
//   }
// }
// `, osNetworkID)
// }

// func testAccComputeInstanceSecgroupMultiUpdate1() string {
// 	return fmt.Sprintf(`
// resource "vkcs_compute_secgroup" "secgroup_1" {
//   name = "secgroup_1"
//   description = "a security group"
//   rule {
//     from_port = 22
//     to_port = 22
//     ip_protocol = "tcp"
//     cidr = "0.0.0.0/0"
//   }
// }

// resource "vkcs_compute_secgroup" "secgroup_2" {
//   name = "secgroup_2"
//   description = "another security group"
//   rule {
//     from_port = 80
//     to_port = 80
//     ip_protocol = "tcp"
//     cidr = "0.0.0.0/0"
//   }
// }

// resource "vkcs_compute_instance" "instance_1" {
//   name = "instance_1"
//   security_groups = ["default"]
//   network {
//     uuid = "%s"
//   }
// }
// `, osNetworkID)
// }

// func testAccComputeInstanceSecgroupMultiUpdate2() string {
// 	return fmt.Sprintf(`
// resource "vkcs_compute_secgroup" "secgroup_1" {
//   name = "secgroup_1"
//   description = "a security group"
//   rule {
//     from_port = 22
//     to_port = 22
//     ip_protocol = "tcp"
//     cidr = "0.0.0.0/0"
//   }
// }

// resource "vkcs_compute_secgroup" "secgroup_2" {
//   name = "secgroup_2"
//   description = "another security group"
//   rule {
//     from_port = 80
//     to_port = 80
//     ip_protocol = "tcp"
//     cidr = "0.0.0.0/0"
//   }
// }

// resource "vkcs_compute_instance" "instance_1" {
//   name = "instance_1"
//   security_groups = ["default", "${vkcs_compute_secgroup.secgroup_1.name}", "${vkcs_compute_secgroup.secgroup_2.name}"]
//   network {
//     uuid = "%s"
//   }
// }
// `, osNetworkID)
// }

func testAccComputeInstanceBootFromVolumeImage() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeInstanceBootFromVolumeVolume() string {
	return fmt.Sprintf(`
resource "vkcs_blockstorage_volume" "vol_1" {
  name = "vol_1"
  size = 5
  image_id = "%s"
  availability_zone = "nova"
  volume_type = "%s"
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "${vkcs_blockstorage_volume.vol_1.id}"
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }

}
`, osImageID, osVolumeType, osNetworkID)
}

func testAccComputeInstanceBootFromVolumeForceNew1() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeInstanceBootFromVolumeForceNew2() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 4
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeInstanceBlockDeviceNewVolume() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
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
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeInstanceBlockDeviceNewVolumeTypeAndBus() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
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
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeInstanceBlockDeviceExistingVolume() string {
	return fmt.Sprintf(`
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "nova"
  volume_type = "%s"
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
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
    uuid = "%s"
  }
}
`, osVolumeType, osImageID, osNetworkID)
}

func testAccComputeInstancePersonality() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
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
    uuid = "%s"
  }
}
`, osNetworkID)
}

// func testAccComputeInstanceMultiEphemeral() string {
// 	return fmt.Sprintf(`
// resource "vkcs_compute_instance" "instance_1" {
//   name = "terraform-test"
//   security_groups = ["default"]
//   block_device {
//     boot_index = 0
//     delete_on_termination = true
//     destination_type = "local"
//     source_type = "image"
//     uuid = "%s"
//   }
//   block_device {
//     boot_index = -1
//     delete_on_termination = true
//     destination_type = "local"
//     source_type = "blank"
//     volume_size = 1
//   }
//   block_device {
//     boot_index = -1
//     delete_on_termination = true
//     destination_type = "local"
//     source_type = "blank"
//     volume_size = 1
//   }
//   network {
//     uuid = "%s"
//   }
// }
// `, osImageID, osNetworkID)
// }

func testAccComputeInstanceAccessIPv4() string {
	return fmt.Sprintf(`
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
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${vkcs_networking_network.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
    access_network = true
  }
}
`, osNetworkID)
}

func testAccComputeInstanceChangeFixedIP1() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.134"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceChangeFixedIP2() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.135"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceStopBeforeDestroy() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  stop_before_destroy = true
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceDetachPortsBeforeDestroy() string {
	return fmt.Sprintf(`

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  vendor_options {
    detach_ports_before_destroy = true
  }
  network {
    port = "${vkcs_networking_port.port_1.id}"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceMetadataRemove1() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    abc = "def"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceMetadataRemove2() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceForceDelete() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  force_delete = true
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceTimeout() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  timeouts {
    create = "10m"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceNetworkModeAuto() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"

  network_mode = "auto"
}
`)
}

func testAccComputeInstanceNetworkModeNone() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "test-instance-1"

  network_mode = "none"
}
`)
}

func testAccComputeInstanceNetworkNameToID() string {
	return fmt.Sprintf(`
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
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    name = "${vkcs_networking_network.network_1.name}"
  }

}
`, osNetworkID)
}

func testAccComputeInstanceCrazyNICs() string {
	return fmt.Sprintf(`
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
    "vkcs_networking_subnet.subnet_1",
    "vkcs_networking_subnet.subnet_2",
    "vkcs_networking_port.port_1",
    "vkcs_networking_port.port_2",
  ]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
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
}
`, osNetworkID)
}

func testAccComputeInstanceStateActive() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "active"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceStateShutoff() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shutoff"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceStateShelve() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shelved_offloaded"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeInstanceTagsCreate() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3"]
}
`, osNetworkID)
}

func testAccComputeInstanceTagsAdd() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3", "tag4"]
}
`, osNetworkID)
}

func testAccComputeInstanceTagsDelete() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag2", "tag3"]
}
`, osNetworkID)
}

func testAccComputeInstanceTagsClear() string {
	return fmt.Sprintf(`
resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}
