package lb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBMember_importBasic(t *testing.T) {
	memberResourceName := "vkcs_lb_member.member_1"
	poolResourceName := "vkcs_lb_pool.pool_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMemberConfigBasic,
			},

			{
				ResourceName:      memberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBMemberImportID(poolResourceName, memberResourceName),
			},
		},
	})
}

func testAccLBMemberImportID(poolResource, memberResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources[poolResource]
		if !ok {
			return "", fmt.Errorf("Pool not found: %s", poolResource)
		}

		member, ok := s.RootModule().Resources[memberResource]
		if !ok {
			return "", fmt.Errorf("Member not found: %s", memberResource)
		}

		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}
