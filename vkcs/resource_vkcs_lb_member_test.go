package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func TestAccLBMember_basic(t *testing.T) {
	var member1 pools.Member
	var member2 pools.Member

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMemberConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMemberExists("vkcs_lb_member.member_1", &member1),
					testAccCheckLBMemberExists("vkcs_lb_member.member_2", &member2),
				),
			},
			{
				Config: TestAccLbMemberConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_lb_member.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("vkcs_lb_member.member_2", "weight", "15"),
				),
			},
		},
	})
}

func testAccCheckLBMemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_member" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(lbClient, poolID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBMemberExists(n string, member *pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(clients.Config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		poolID := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(lbClient, poolID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*member = *found

		return nil
	}
}

const TestAccLbMemberConfigBasic = `
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

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
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

resource "vkcs_lb_member" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  pool_id = vkcs_lb_pool.pool_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id
  weight = 1

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "vkcs_lb_member" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  pool_id = vkcs_lb_pool.pool_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbMemberConfigUpdate = `
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

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
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

resource "vkcs_lb_member" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  weight = 10
  admin_state_up = "true"
  pool_id = vkcs_lb_pool.pool_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "vkcs_lb_member" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  weight = 15
  admin_state_up = "true"
  pool_id = vkcs_lb_pool.pool_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
