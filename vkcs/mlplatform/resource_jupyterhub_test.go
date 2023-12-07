package mlplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccMLPlatformJupyterHub_resize_big(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccMLPlatformJupyterHubResize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "name", "tfacc-jupyter-hub"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "boot_volume.name", "ml_platform_boot_volume"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "boot_volume.size", "50"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.#", "2"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.0.name", "ml_platform_data_volume"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.1.name", "ml_platform_volume_2"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.1.size", "70"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccMLPlatformJupyterHubResizeUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "name", "tfacc-jupyter-hub"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "boot_volume.name", "ml_platform_boot_volume"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "boot_volume.size", "51"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.#", "2"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.0.name", "ml_platform_data_volume"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.0.size", "61"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.1.name", "ml_platform_volume_2"),
					resource.TestCheckResourceAttr("vkcs_mlplatform_jupyterhub.jupyterhub", "data_volumes.1.size", "71"),
				),
			},
			{
				ResourceName:            "vkcs_mlplatform_jupyterhub.jupyterhub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password", "networks"},
			},
		},
	})
}

const testAccMLPlatformJupyterHubResize = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_mlplatform_jupyterhub" "jupyterhub" {
    name = "tfacc-jupyter-hub"
    domain_name = "mlp-tf-example"
	admin_name = "admin"
    admin_password = "dM8Ao21,0S264iZp"
    flavor_id = data.vkcs_compute_flavor.base.id
    availability_zone = "{{.AvailabilityZone}}"
    boot_volume = {
        size = 50
        volume_type = "ceph-ssd"
    }
    data_volumes = [
        {
            size = 60
            volume_type = "ceph-ssd"
        },
        {
            size = 70
            volume_type = "ceph-ssd"
        }
    ]
    networks = [
        {
            network_id= vkcs_networking_network.base.id
        },
    ]
}
`

const testAccMLPlatformJupyterHubResizeUpdate = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_mlplatform_jupyterhub" "jupyterhub" {
    name = "tfacc-jupyter-hub"
    admin_name = "admin"
    admin_password = "dM8Ao21,0S264iZp"
    flavor_id = data.vkcs_compute_flavor.base.id
    availability_zone = "{{.AvailabilityZone}}"
    boot_volume = {
        size = 51
        volume_type = "ceph-ssd"
    }
    data_volumes = [
        {
            size = 61
            volume_type = "ceph-ssd"
        },
        {
            size = 71
            volume_type = "ceph-ssd"
        }
    ]
    networks = [
        {
            network_id= vkcs_networking_network.base.id
        },
    ]
}
`
