package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeInterfaceAttachImport_basic(t *testing.T) {
	resourceName := "vkcs_compute_interface_attach.ai_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInterfaceAttachBasic, testAccValues),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_pass",
				},
			},
		},
	})
}
