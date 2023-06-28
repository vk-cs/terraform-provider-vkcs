package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBMembers_importBasic(t *testing.T) {
	membersResourceName := "vkcs_lb_members.members_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMembersConfigBasic,
			},

			{
				ResourceName:      membersResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
