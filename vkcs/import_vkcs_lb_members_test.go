package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBMembers_importBasic(t *testing.T) {
	membersResourceName := "vkcs_lb_members.members_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
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
