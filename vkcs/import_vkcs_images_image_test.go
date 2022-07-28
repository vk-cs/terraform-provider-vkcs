package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImagesImage_importBasic(t *testing.T) {
	resourceName := "vkcs_images_image.image_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"local_file_path",
					"image_cache_path",
					"image_source_url",
					"verify_checksum",
				},
			},
		},
	})
}
