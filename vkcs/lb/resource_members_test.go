package lb_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccLBMembers_basic(t *testing.T) {
	var members []pools.Member

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLBMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMembersConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.110",
						"weight":  "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.111",
						"weight":  "1",
						"backup":  "true",
					}),
				),
			},
			{
				Config: TestAccLbMembersConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.110",
						"weight":  "10",
						"backup":  "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.111",
						"weight":  "15",
						"backup":  "false",
					}),
				),
			},
			{
				Config: TestAccLbMembersConfigUnsetSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.110",
						"weight":  "10",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address": "192.168.199.111",
						"weight":  "15",
					}),
				),
			},
			{
				Config: TestAccLbMembersConfigDeleteMembers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "0"),
				),
			},
		},
	})
}

func TestAccLBMembers_replaceInstance(t *testing.T) {
	var members []pools.Member

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLBMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLbMembersConfigInstances, map[string]string{"UserData": "rev-1"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address":       "192.168.199.20",
						"protocol_port": "8080",
					}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLbMembersConfigInstances, map[string]string{"UserData": "rev-2"}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_compute_instance.instance[0]", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("vkcs_lb_members.members_1", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address":       "192.168.199.20",
						"protocol_port": "8080",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vkcs_lb_members.members_1", "member.*", map[string]string{
						"address":       "192.168.199.21",
						"protocol_port": "8080",
					}),
				),
			},
		},
	})
}

func TestAccLBMembers_migrateToFramework(t *testing.T) {
	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCheckLBMembersDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"vkcs": {
						VersionConstraint: "0.19.0",
						Source:            "vk-cs/vkcs",
					},
				},
				Config: TestAccLbMembersConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   TestAccLbMembersConfigBasic,
				PlanOnly:                 true,
			},
		},
	})
}

func testAccLBMembersClient() (clients.Config, error) {
	opts := clients.ConfigOpts{}
	return opts.LoadAndValidate()
}

func testAccCheckLBMembersDestroy(s *terraform.State) error {
	config, err := testAccLBMembersClient()
	if err != nil {
		return err
	}
	lbClient, err := config.LoadBalancerV2Client(config.GetRegion())
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_members" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]

		allPages, err := pools.ListMembers(lbClient, poolID, pools.ListMembersOpts{}).AllPages()
		if err != nil {
			if errutil.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("Error getting vkcs_lb_members: %s", err)
		}

		members, err := pools.ExtractMembers(allPages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve vkcs_lb_members: %s", err)
		}

		if len(members) > 0 {
			return fmt.Errorf("Members still exist: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBMembersExists(n string, members *[]pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config, err := testAccLBMembersClient()
		if err != nil {
			return err
		}
		lbClient, err := config.LoadBalancerV2Client(config.GetRegion())
		if err != nil {
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		poolID := rs.Primary.Attributes["pool_id"]
		allPages, err := pools.ListMembers(lbClient, poolID, pools.ListMembersOpts{}).AllPages()
		if err != nil {
			return fmt.Errorf("Error getting vkcs_lb_members: %s", err)
		}

		found, err := pools.ExtractMembers(allPages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve vkcs_lb_members: %s", err)
		}

		*members = found

		return nil
	}
}

const TestAccLbMembersConfigBasic = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = vkcs_networking_network.network_1.id
  cidr = "192.168.199.0/24"

}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id
  vip_address = "192.168.199.10"
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_members" "members_1" {
  pool_id = vkcs_lb_pool.pool_1.id

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    subnet_id = vkcs_networking_subnet.subnet_1.id
    weight = 0
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
	subnet_id = vkcs_networking_subnet.subnet_1.id
	backup = true
  }

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbMembersConfigUpdate = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_members" "members_1" {
  pool_id = vkcs_lb_pool.pool_1.id

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    weight = 10
    admin_state_up = "true"
    subnet_id = vkcs_networking_subnet.subnet_1.id
	backup = true
}

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
	subnet_id = vkcs_networking_subnet.subnet_1.id
	backup = false
  }

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbMembersConfigUnsetSubnet = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_members" "members_1" {
  pool_id = vkcs_lb_pool.pool_1.id

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    weight = 10
    admin_state_up = "true"
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
  }

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbMembersConfigDeleteMembers = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_members" "members_1" {
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const testAccLbMembersConfigInstances = `
data "vkcs_networking_network" "extnet" {
  name = "{{.ExtNetName}}"
}

resource "vkcs_networking_network" "base" {
  name           = "lb-members-net"
  admin_state_up = true
  sdn            = data.vkcs_networking_network.extnet.sdn
}

resource "vkcs_networking_subnet" "base" {
  name       = "lb-members-subnet"
  network_id = vkcs_networking_network.base.id
  cidr       = "192.168.199.0/24"
  sdn        = data.vkcs_networking_network.extnet.sdn
}

resource "vkcs_networking_router" "base" {
  name                = "lb-members-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
  sdn                 = data.vkcs_networking_network.extnet.sdn
}

resource "vkcs_networking_router_interface" "base" {
  router_id = vkcs_networking_router.base.id
  subnet_id = vkcs_networking_subnet.base.id
  sdn       = data.vkcs_networking_network.extnet.sdn
}

{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance" {
  count              = 1
  depends_on         = [vkcs_networking_router_interface.base]
  name               = "lb-member-${count.index}"
  availability_zone  = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  flavor_id          = data.vkcs_compute_flavor.base.id
  user_data          = count.index == 0 ? "{{.UserData}}" : "static"

  block_device {
    uuid                  = data.vkcs_images_image.base.id
    source_type           = "image"
    destination_type      = "volume"
    volume_type           = "{{.VolumeType}}"
    volume_size           = 10
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    uuid        = vkcs_networking_network.base.id
    fixed_ip_v4 = "192.168.199.2${count.index}"
  }
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.base.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_members" "members_1" {
  pool_id = vkcs_lb_pool.pool_1.id

  dynamic "member" {
    for_each = vkcs_compute_instance.instance
    content {
      address       = member.value.access_ip_v4
      protocol_port = 8080
    }
  }

  member {
    address       = "192.168.199.21"
    protocol_port = 8080
  }
}
`
