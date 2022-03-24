package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func testAccCheckLBMembersComputeHash(members *[]pools.Member, weight int, address string, idx *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		membersResource := resourceMembers().Schema["member"].Elem.(*schema.Resource)
		f := schema.HashResource(membersResource)

		for _, m := range flattenLBMembers(*members) {
			if m["address"] == address && m["weight"] == weight {
				*idx = f(m)
				break
			}
		}

		return nil
	}
}

// func testCheckResourceAttrWithIndexesAddr(name, format string, idx *int, value string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		return resource.TestCheckResourceAttr(name, fmt.Sprintf(format, *idx), value)(s)
// 	}
// }

// func testCheckResourceAttrSetWithIndexesAddr(name, format string, idx *int) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		return resource.TestCheckResourceAttrSet(name, fmt.Sprintf(format, *idx))(s)
// 	}
// }

func TestAccLBMembers_basic(t *testing.T) {
	var members []pools.Member
	var idx1 int
	var idx2 int

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMembersConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					testAccCheckLBMembersComputeHash(&members, 0, "192.168.199.110", &idx1),
					testAccCheckLBMembersComputeHash(&members, 1, "192.168.199.111", &idx2),
					// checks for TypeSet are currently not supported by SDK2
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx1, "0"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx2, "1"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.backup", &idx1, "false"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.backup", &idx2, "true"),
					// testCheckResourceAttrSetWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx1),
					// testCheckResourceAttrSetWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx2),
				),
			},
			{
				Config: TestAccLbMembersConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					testAccCheckLBMembersComputeHash(&members, 10, "192.168.199.110", &idx1),
					testAccCheckLBMembersComputeHash(&members, 15, "192.168.199.111", &idx2),
					// checks for TypeSet are currently not supported by SDK2
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx1, "10"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx2, "15"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.backup", &idx1, "true"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.backup", &idx2, "false"),
					// testCheckResourceAttrSetWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx1),
					// testCheckResourceAttrSetWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx2),
				),
			},
			{
				Config: TestAccLbMembersConfigUnsetSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMembersExists("vkcs_lb_members.members_1", &members),
					resource.TestCheckResourceAttr("vkcs_lb_members.members_1", "member.#", "2"),
					testAccCheckLBMembersComputeHash(&members, 10, "192.168.199.110", &idx1),
					testAccCheckLBMembersComputeHash(&members, 15, "192.168.199.111", &idx2),
					// checks for TypeSet are currently not supported by SDK2
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx1, "10"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.weight", &idx2, "15"),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx1, ""),
					// testCheckResourceAttrWithIndexesAddr("vkcs_lb_members.members_1", "member.%d.subnet_id", &idx2, ""),
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

func testAccCheckLBMembersDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_members" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]

		allPages, err := pools.ListMembers(lbClient, poolID, pools.ListMembersOpts{}).AllPages()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
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

		config := testAccProvider.Meta().(*config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
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
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
  vip_address = "192.168.199.10"
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener_1.id}"
}

resource "vkcs_lb_members" "members_1" {
  pool_id = "${vkcs_lb_pool.pool_1.id}"

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    weight = 0
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
	subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
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
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener_1.id}"
}

resource "vkcs_lb_members" "members_1" {
  pool_id = "${vkcs_lb_pool.pool_1.id}"

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    weight = 10
    admin_state_up = "true"
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
	backup = true
}

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
	subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
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
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener_1.id}"
}

resource "vkcs_lb_members" "members_1" {
  pool_id = "${vkcs_lb_pool.pool_1.id}"

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
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener_1.id}"
}

resource "vkcs_lb_members" "members_1" {
  pool_id = "${vkcs_lb_pool.pool_1.id}"

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`
