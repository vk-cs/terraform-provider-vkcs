package iam_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/serviceusers"
)

func TestAccIAMServiceUser_basic(t *testing.T) {
	var serviceUser serviceusers.ServiceUser

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMServiceUserBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceUserExists("vkcs_iam_service_user.basic", &serviceUser),
					resource.TestCheckResourceAttr("vkcs_iam_service_user.basic", "name", "tfacc-service-user-basic"),
					resource.TestCheckResourceAttr("vkcs_iam_service_user.basic", "description", "Service user created by acceptance test"),
					acctest.TestCheckResourceListAttr("vkcs_iam_service_user.basic", "role_names", []string{"mcs_admin_vm", "mcs_admin_network"}),
					resource.TestCheckResourceAttrSet("vkcs_iam_service_user.basic", "created_at"),
					resource.TestCheckResourceAttrPtr("vkcs_iam_service_user.basic", "creator_name", &serviceUser.CreatorName),
					resource.TestCheckResourceAttrSet("vkcs_iam_service_user.basic", "login"),
					resource.TestCheckResourceAttrSet("vkcs_iam_service_user.basic", "password"),
				),
			},
			acctest.ImportStep("vkcs_iam_service_user.basic", "login", "password"),
		},
	})
}

func testAccCheckServiceUserExists(n string, serviceUser *serviceusers.ServiceUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		config, err := clients.ConfigureFromEnv(context.Background())
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.IAMServiceUsersV1Client(config.GetRegion())
		if err != nil {
			return fmt.Errorf("Error creating VKCS IAM Service Users client: %s", err)
		}

		found, err := serviceusers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*serviceUser = *found

		return nil
	}
}

const testAccIAMServiceUserBasic = `
resource "vkcs_iam_service_user" "basic" {
  name        = "tfacc-service-user-basic"
  description = "Service user created by acceptance test"
  role_names  = ["mcs_admin_vm", "mcs_admin_network"]
}
`
