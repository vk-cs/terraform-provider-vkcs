package backup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBackupProviderDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccBackupProviderDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_backup_provider.provider", "name", "cloud_servers"),
				),
			},
		},
	})
}

const testAccBackupProviderDataSourceBasic = `
data "vkcs_backup_provider" "provider" {
	name = "cloud_servers"
}
`
