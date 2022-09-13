package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeyManagerContainer_importBasic(t *testing.T) {
	resourceName := "vkcs_keymanager_container.container_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccKeyManagerContainerBasic, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer}),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyManagerContainer_importACLs(t *testing.T) {
	resourceName := "vkcs_keymanager_container.container_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccKeyManagerContainerAcls, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer}),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
