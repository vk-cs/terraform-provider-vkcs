package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func TestAccLBPool_basic(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckLB(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbPoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBPoolExists("vkcs_lb_pool.pool_1", &pool),
				),
			},
			{
				Config: TestAccLbPoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_lb_pool.pool_1", "name", "pool_1_updated"),
				),
			},
		},
	})
}

func TestAccLBPool_octavia_udp(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbPoolConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBPoolExists("vkcs_lb_pool.pool_1", &pool),
					resource.TestCheckResourceAttr("vkcs_lb_pool.pool_1", "protocol", "UDP"),
				),
			},
		},
	})
}

func testAccCheckLBPoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_pool" {
			continue
		}

		_, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBPoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
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

		found, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*pool = *found

		return nil
	}
}

const TestAccLbPoolConfigBasic = `
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
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbPoolConfigUpdate = `
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
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1_updated"
  protocol = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id = "${vkcs_lb_listener.listener_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbPoolConfigOctaviaUDP = `
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

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "UDP"
  lb_method = "ROUND_ROBIN"
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
