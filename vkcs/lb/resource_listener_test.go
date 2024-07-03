package lb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ilisteners "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/lb/v2/listeners"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
)

func TestAccLBListener_basic(t *testing.T) {
	var listener listeners.Listener
	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerBasic, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "connection_limit", "-1"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_fullUpdate(t *testing.T) {
	var listener listeners.Listener
	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerFullUpdateOld, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerFullUpdateNew, map[string]string{"TestAccLBListenerBase": baseConfig}),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_octavia(t *testing.T) {
	var listener listeners.Listener
	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctavia, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "connection_limit", "5"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_client_data", "1000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_member_connect", "2000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_member_data", "3000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_tcp_inspect", "4000"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctaviaUpdate, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "connection_limit", "100"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_client_data", "4000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_member_connect", "3000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_member_data", "2000"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "timeout_tcp_inspect", "1000"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_octaviaUDP(t *testing.T) {
	var listener listeners.Listener
	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctaviaUDP, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "protocol", "UDP"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_octaviaInsertHeaders(t *testing.T) {
	var listener listeners.Listener
	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctaviaInsertHeaders1, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-For", "true"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-Port", "false"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctaviaInsertHeaders2, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-For", "false"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-Port", "true"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerOctavia, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckNoResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-For"),
					resource.TestCheckNoResourceAttr(
						"vkcs_lb_listener.listener_1", "insert_headers.X-Forwarded-Port"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_allowedCidrsOrderIgnored(t *testing.T) {
	var listener listeners.Listener

	baseConfig := acctest.AccTestRenderConfig(testAccLBListenerBase)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerAllowedCidrsOrderIgnoredOld, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBListenerExists("vkcs_lb_listener.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "allowed_cidrs.0", "192.168.1.0/24"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccLBListenerAllowedCidrsOrderIgnoredNew, map[string]string{"TestAccLBListenerBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_lb_listener.listener_1", "allowed_cidrs.0", "192.168.1.0/24"),
				),
			},
			{
				ResourceName:      "vkcs_lb_listener.listener_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLBListenerDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_listener" {
			continue
		}

		_, err := ilisteners.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		lbClient, err := config.LoadBalancerV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		found, err := ilisteners.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

const testAccLBListenerBase = `
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLBListenerBasic = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLBListenerFullUpdateOld = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name                   = "listener_1"
  protocol               = "HTTP"
  protocol_port          = 8080
  connection_limit       = 100
  admin_state_up         = "true"
  loadbalancer_id        = vkcs_lb_loadbalancer.loadbalancer_1.id
  default_pool_id        = vkcs_lb_pool.pool_1.id
  description            = "My old listener"
  timeout_client_data    = 1000
  timeout_member_connect = 2000
  timeout_member_data    = 3000
  timeout_tcp_inspect    = 4000
  insert_headers = {
    X-Forwarded-For = "true"
  }
  allowed_cidrs = ["192.168.1.0/24"]
}

`

const testAccLBListenerFullUpdateNew = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_pool" "pool_2" {
  name            = "pool_2"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name                   = "listener_2"
  protocol               = "HTTP"
  protocol_port          = 8080
  loadbalancer_id        = vkcs_lb_loadbalancer.loadbalancer_1.id
  connection_limit       = 200
  admin_state_up         = "false"
  default_pool_id        = vkcs_lb_pool.pool_2.id
  description            = "My new listener"
  timeout_client_data    = 2000
  timeout_member_connect = 3000
  timeout_member_data    = 4000
  timeout_tcp_inspect    = 5000
  insert_headers         = {}
  allowed_cidrs          = ["192.168.2.0/24"]
}
`

const testAccLBListenerOctavia = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name                   = "listener_1"
  protocol               = "HTTP"
  protocol_port          = 8080
  connection_limit       = 5
  timeout_client_data    = 1000
  timeout_member_connect = 2000
  timeout_member_data    = 3000
  timeout_tcp_inspect    = 4000
  default_pool_id        = vkcs_lb_pool.pool_1.id
  loadbalancer_id        = vkcs_lb_loadbalancer.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLBListenerOctaviaUpdate = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name                   = "listener_1_updated"
  protocol               = "HTTP"
  protocol_port          = 8080
  connection_limit       = 100
  timeout_client_data    = 4000
  timeout_member_connect = 3000
  timeout_member_data    = 2000
  timeout_tcp_inspect    = 1000
  admin_state_up         = "true"
  default_pool_id        = vkcs_lb_pool.pool_1.id
  loadbalancer_id        = vkcs_lb_loadbalancer.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLBListenerOctaviaUDP = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "UDP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "UDP"
  protocol_port   = 53
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLBListenerOctaviaInsertHeaders1 = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id

  insert_headers = {
    X-Forwarded-For  = "true"
    X-Forwarded-Port = "false"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLBListenerOctaviaInsertHeaders2 = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id

  insert_headers = {
    X-Forwarded-For  = "false"
    X-Forwarded-Port = "true"
  }
}
`

const testAccLBListenerAllowedCidrsOrderIgnoredOld = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
  allowed_cidrs   = ["192.168.1.0/24", "192.168.2.0/24"]
}
`

const testAccLBListenerAllowedCidrsOrderIgnoredNew = `
{{ .TestAccLBListenerBase }}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  default_pool_id = vkcs_lb_pool.pool_1.id
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
  allowed_cidrs   = ["192.168.2.0/24", "192.168.1.0/24"]
}
`
