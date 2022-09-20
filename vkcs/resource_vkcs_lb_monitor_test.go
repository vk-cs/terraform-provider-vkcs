package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
)

func TestAccLBMonitor_basic(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMonitorConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMonitorExists(t, "vkcs_lb_monitor.monitor_1", &monitor),
				),
			},
			{
				Config: TestAccLbMonitorConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_lb_monitor.monitor_1", "name", "monitor_1_updated"),
					resource.TestCheckResourceAttr("vkcs_lb_monitor.monitor_1", "delay", "30"),
					resource.TestCheckResourceAttr("vkcs_lb_monitor.monitor_1", "timeout", "15"),
				),
			},
		},
	})
}

func TestAccLBMonitor_octavia(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMonitorConfigOctavia,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMonitorExists(t, "vkcs_lb_monitor.monitor_1", &monitor),
					resource.TestCheckResourceAttr("vkcs_lb_monitor.monitor_1", "max_retries_down", "8"),
				),
			},
			{
				Config: TestAccLbMonitorConfigOctaviaUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_lb_monitor.monitor_1", "name", "monitor_1_updated"),
					resource.TestCheckResourceAttr("vkcs_lb_monitor.monitor_1", "max_retries_down", "3"),
				),
			},
		},
	})
}

func TestAccLBMonitor_octavia_udp(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMonitorConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBMonitorExists(t, "vkcs_lb_monitor.monitor_1", &monitor),
					resource.TestCheckResourceAttr("vkcs_lb_monitor.monitor_1", "type", "UDP-CONNECT"),
				),
			},
		},
	})
}

func testAccCheckLBMonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_monitor" {
			continue
		}

		_, err := monitors.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBMonitorExists(t *testing.T, n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		found, err := monitors.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

const TestAccLbMonitorConfigBasic = `
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

resource "vkcs_lb_monitor" "monitor_1" {
  name = "monitor_1"
  type = "PING"
  delay = 20
  timeout = 10
  max_retries = 5
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbMonitorConfigUpdate = `
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

resource "vkcs_lb_monitor" "monitor_1" {
  name = "monitor_1_updated"
  type = "PING"
  delay = 30
  timeout = 15
  max_retries = 10
  admin_state_up = "true"
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbMonitorConfigOctavia = `
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

resource "vkcs_lb_monitor" "monitor_1" {
  name = "monitor_1"
  type = "PING"
  delay = 20
  timeout = 10
  max_retries = 5
  max_retries_down = 8
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbMonitorConfigOctaviaUpdate = `
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

resource "vkcs_lb_monitor" "monitor_1" {
  name = "monitor_1_updated"
  type = "PING"
  delay = 30
  timeout = 15
  max_retries = 10
  max_retries_down = 3
  admin_state_up = "true"
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbMonitorConfigOctaviaUDP = `
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
  protocol = "UDP"
  protocol_port = 53
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_1" {
  name = "pool_1"
  protocol = "UDP"
  lb_method = "ROUND_ROBIN"
  listener_id = vkcs_lb_listener.listener_1.id
}

resource "vkcs_lb_monitor" "monitor_1" {
  name = "monitor_1"
  type = "UDP-CONNECT"
  delay = 20
  timeout = 10
  max_retries = 5
  pool_id = vkcs_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
