package vkcs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccLBLoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbLoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBLoadBalancerExists("vkcs_lb_loadbalancer.loadbalancer_1", &lb),
					testAccCheckLBLoadBalancerHasTag("vkcs_lb_loadbalancer.loadbalancer_1", "tag1"),
					testAccCheckLBLoadBalancerTagCount("vkcs_lb_loadbalancer.loadbalancer_1", 1),
				),
			},
			{
				Config: testAccLbLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBLoadBalancerHasTag("vkcs_lb_loadbalancer.loadbalancer_1", "tag1"),
					testAccCheckLBLoadBalancerHasTag("vkcs_lb_loadbalancer.loadbalancer_1", "tag2"),
					testAccCheckLBLoadBalancerTagCount("vkcs_lb_loadbalancer.loadbalancer_1", 2),
					resource.TestCheckResourceAttr("vkcs_lb_loadbalancer.loadbalancer_1", "name", "loadbalancer_1_updated"),
					resource.TestMatchResourceAttr("vkcs_lb_loadbalancer.loadbalancer_1", "vip_port_id", regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
		},
	})
}

func TestAccLBLoadBalancer_vip_network(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbLoadBalancerConfigVIPNetwork,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBLoadBalancerExists("vkcs_lb_loadbalancer.loadbalancer_1", &lb),
				),
			},
		},
	})
}

func TestAccLBLoadBalancer_vip_port_id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var port ports.Port

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbLoadBalancerConfigVIPPortID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBLoadBalancerExists("vkcs_lb_loadbalancer.loadbalancer_1", &lb),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttrPtr("vkcs_lb_loadbalancer.loadbalancer_1", "vip_port_id", &port.ID),
				),
			},
		},
	})
}

func testAccCheckLBLoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_loadbalancer" {
			continue
		}

		lb, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil && lb.ProvisioningStatus != "DELETED" {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBLoadBalancerExists(
	n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		*lb = *found

		return nil
	}
}
func testAccCheckLBLoadBalancerHasTag(n, tag string) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckLBLoadBalancerTagCount(n string, expected int) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

const testAccLbLoadBalancerConfigBasic = `
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
	  tags = ["tag1"]

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`

const testAccLbLoadBalancerConfigUpdate = `
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
      name = "loadbalancer_1_updated"
      admin_state_up = "true"
      vip_subnet_id = vkcs_networking_subnet.subnet_1.id
	  tags = ["tag1", "tag2"]

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`

const testAccLbLoadBalancerConfigVIPNetwork = `
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
  vip_network_id = vkcs_networking_network.network_1.id
  depends_on = ["vkcs_networking_subnet.subnet_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbLoadBalancerConfigVIPPortID = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_port" "port_1" {
  name           = "port_1"
  network_id     = vkcs_networking_network.network_1.id
  admin_state_up = "true"
  depends_on = ["vkcs_networking_subnet.subnet_1"]
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_port_id = vkcs_networking_port.port_1.id
  depends_on = ["vkcs_networking_port.port_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`
