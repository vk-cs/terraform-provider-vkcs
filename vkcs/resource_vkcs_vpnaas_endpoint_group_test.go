package vkcs

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func TestAccVPNaaSGroup_basic(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupExists("vkcs_vpnaas_endpoint_group.group_1", &group),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "type", &group.Type),
					testAccCheckEndpoints("vkcs_vpnaas_endpoint_group.group_1", &group.Endpoints),
				),
			},
		},
	})
}

func TestAccVPNaaSGroup_update(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupExists("vkcs_vpnaas_endpoint_group.group_1", &group),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "type", &group.Type),
					testAccCheckEndpoints("vkcs_vpnaas_endpoint_group.group_1", &group.Endpoints),
				),
			},
			{
				Config: testAccEndpointGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupExists("vkcs_vpnaas_endpoint_group.group_1", &group),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_endpoint_group.group_1", "type", &group.Type),
					testAccCheckEndpoints("vkcs_vpnaas_endpoint_group.group_1", &group.Endpoints),
				),
			},
		},
	})
}

func testAccCheckEndpointGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_group" {
			continue
		}
		_, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EndpointGroup (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckEndpointGroupExists(n string, group *endpointgroups.EndpointGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var found *endpointgroups.EndpointGroup

		found, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*group = *found

		return nil
	}
}

func testAccCheckEndpoints(n string, actual *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		var endpointsList []string
		// Find all "endpoints.<number>" keys and collect the values.
		// The <number> values are seemingly random and very large.
		for k, v := range rs.Primary.Attributes {
			println("[DEBUG] key:", k, "value:", v)
			if strings.HasPrefix(k, "endpoints.") && k[10] >= '0' && k[10] <= '9' {
				endpointsList = append(endpointsList, v)
			}
		}

		if len(*actual) != len(endpointsList) {
			return fmt.Errorf("The number of endpoints did not match. Expected: %v but got %v", len(*actual), len(endpointsList))
		}

		sort.Strings(endpointsList)
		sort.Strings(*actual)

		for i, endpoint := range endpointsList {
			if endpoint != (*actual)[i] {
				return fmt.Errorf("The endpoints did not match. Expected: '%v' but got '%v'", endpoint, (*actual)[i])
			}
		}
		return nil
	}
}

const testAccEndpointGroupBasic = `
	resource "vkcs_vpnaas_endpoint_group" "group_1" {
		name = "Group 1"
		type = "cidr"
		endpoints = ["10.3.0.0/24",
			"10.2.0.0/24",]
	}
`

const testAccEndpointGroupUpdate = `
	resource "vkcs_vpnaas_endpoint_group" "group_1" {
		name = "Updated Group 1"
		type = "cidr"
		endpoints = ["10.2.0.0/24",
			"10.3.0.0/24",]
	}
`
