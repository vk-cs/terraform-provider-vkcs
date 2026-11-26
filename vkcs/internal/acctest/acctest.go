package acctest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
)

var (
	OsFlavorName        = os.Getenv("OS_FLAVOR_NAME")
	OsNewFlavorName     = os.Getenv("OS_NEW_FLAVOR_NAME")
	OsImageName         = os.Getenv("OS_IMAGE_NAME")
	OsRegionName        = os.Getenv("OS_REGION_NAME")
	OsProjectID         = os.Getenv("OS_PROJECT_ID")
	OsExtNetName        = os.Getenv("OS_EXT_NET_NAME")
	OsExtNetNameNeutron = os.Getenv("OS_EXT_NET_NAME_NEUTRON")
	OsAvailabilityZone  = os.Getenv("OS_AVAILABILITY_ZONE")
	OsVolumeType        = os.Getenv("OS_VOLUME_TYPE")
	OsDbVolumeType      = util.GetFirstNotEmpty(os.Getenv("OS_DB_VOLUME_TYPE"), OsVolumeType)
	accTestRunID        = strconv.FormatInt(time.Now().UnixNano(), 36)
)

var AccTestValues map[string]string = map[string]string{
	"BaseNetwork":           AccTestBaseNetwork(),
	"BaseExtNetwork":        AccTestBaseExtNetwork(),
	"BaseExtNetworkNeutron": AccTestBaseExtNetworkNeutron(),
	"BaseImage":             AccTestBaseImage(),
	"BaseFlavor":            AccTestBaseFlavor(),
	"BaseNewFlavor":         AccTestBaseNewFlavor(),
	"BaseSecurityGroup":     AccTestBaseDefaultSG(),
	"BaseBootVolume":        AccTestBaseBootVolume(),
	"AvailabilityZone":      OsAvailabilityZone,
	"VolumeType":            OsVolumeType,
	"DbVolumeType":          OsDbVolumeType,
	"FlavorName":            OsFlavorName,
	"NewFlavorName":         OsNewFlavorName,
	"ImageName":             OsImageName,
	"ExtNetName":            OsExtNetName,
	"ExtNetNameNeutron":     OsExtNetNameNeutron,
	"ProjectID":             OsProjectID,
}

var AccTestProviders map[string]func() (*schema.Provider, error)
var AccTestProvider *schema.Provider
var AccTestProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

func init() {
	AccTestProvider = provider.SDKProvider()
	AccTestProviders = map[string]func() (*schema.Provider, error){
		"vkcs": func() (*schema.Provider, error) {
			return AccTestProvider, nil
		},
	}
	AccTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"vkcs": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(provider.Provider()),
				func() tfprotov6.ProviderServer {
					server, _ := tf5to6server.UpgradeServer(
						ctx,
						provider.SDKProvider().GRPCProvider,
					)
					return server
				},
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer, nil
		},
	}
}

func AccTestPreCheck(t *testing.T) {
	vars := map[string]interface{}{
		"OS_VOLUME_TYPE":          OsVolumeType,
		"OS_AVAILABILITY_ZONE":    OsAvailabilityZone,
		"OS_FLAVOR_NAME":          OsFlavorName,
		"OS_NEW_FLAVOR_NAME":      OsNewFlavorName,
		"OS_IMAGE_NAME":           OsImageName,
		"OS_EXT_NET_NAME":         OsExtNetName,
		"OS_EXT_NET_NAME_NEUTRON": OsExtNetNameNeutron,
		"OS_PROJECT_ID":           OsProjectID,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
}

func AccTestBaseExtNetwork() string {
	return fmt.Sprintf(`
	data "vkcs_networking_network" "extnet" {
		name = "%s"
	  }
	`, OsExtNetName)
}

func AccTestBaseExtNetworkNeutron() string {
	return fmt.Sprintf(`
	data "vkcs_networking_network" "extnet" {
		name = "%s"
	  }
	`, OsExtNetNameNeutron)
}

func AccTestBaseFlavor() string {
	return fmt.Sprintf(`
	data "vkcs_compute_flavor" "base" {
		name = "%s"
	}
`, OsFlavorName)
}

func AccTestBaseNewFlavor() string {
	return fmt.Sprintf(`
	data "vkcs_compute_flavor" "base" {
		name = "%s"
	}
`, OsNewFlavorName)
}

func AccTestBaseBootVolume() string {
	return fmt.Sprintf(`block_device {
    uuid                  = data.vkcs_images_image.base.id
    source_type           = "image"
    destination_type      = "volume"
    volume_type           = "%s"
    volume_size           = 10
    boot_index            = 0
    delete_on_termination = true
  }`, OsVolumeType)
}

func AccTestBaseImage() string {
	return fmt.Sprintf(`
	data "vkcs_images_image" "base" {
		name = "%s"
	}
`, OsImageName)
}

func AccTestBaseNetwork() string {
	return fmt.Sprintf(`

data "vkcs_networking_network" "extnet" {
	name = "%[1]s"
  }
  
  resource "vkcs_networking_network" "base" {
	name           = "base-net-%[2]s"
	admin_state_up = true
	sdn            = data.vkcs_networking_network.extnet.sdn
  }
  
  resource "vkcs_networking_subnet" "base" {
	name       = "subnet_1"
	network_id = vkcs_networking_network.base.id
	cidr       = "192.168.199.0/24"
	sdn        = data.vkcs_networking_network.extnet.sdn
  }
  
  resource "vkcs_networking_router" "base" {
	depends_on          = [vkcs_networking_subnet.base]
	name                = "base-router"
	admin_state_up      = true
	external_network_id = data.vkcs_networking_network.extnet.id
	sdn                 = data.vkcs_networking_network.extnet.sdn
  }
  
  resource "vkcs_networking_router_interface" "base" {
	router_id = vkcs_networking_router.base.id
	subnet_id = vkcs_networking_subnet.base.id
	sdn       = data.vkcs_networking_network.extnet.sdn
  }
`, OsExtNetName, accTestRunID)
}

func AccTestBaseDefaultSG() string {
	return fmt.Sprintf(`
	data "vkcs_networking_secgroup" "default_secgroup" {
	  name = "default"
	  sdn  = "%s"
	}
	`, networking.DefaultSDN)
}

func AccTestRenderConfig(testConfig string, values ...map[string]string) string {
	t := template.Must(template.New("acc").Option("missingkey=error").Parse(testConfig))
	buf := &bytes.Buffer{}

	tmplValues := map[string]string{}
	util.CopyToMap(&tmplValues, &AccTestValues)
	if len(values) > 0 {
		util.CopyToMap(&tmplValues, &values[0])
	}

	_ = t.Execute(buf, tmplValues)

	return buf.String()
}

func AccTestGetStepsWithMigrationCases(steps []resource.TestStep) (migrationSteps []resource.TestStep) {
	for _, step := range steps {
		migrationSteps = append(migrationSteps, resource.TestStep{
			Config: step.Config,
			Check:  step.Check,
			ExternalProviders: map[string]resource.ExternalProvider{
				"vkcs": {
					VersionConstraint: "0.2.1",
					Source:            "vk-cs/vkcs",
				},
			},
		})
	}
	return append(migrationSteps, steps...)
}

func ImportStep(resourceName string, verifyIgnore ...string) resource.TestStep {
	return resource.TestStep{
		ResourceName:            resourceName,
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: verifyIgnore,
	}
}

func GenerateUniqueTestFields(testName string) map[string]string {
	t := time.Now()
	return map[string]string{
		"TestName":    testName,
		"CurrentTime": fmt.Sprintf("%dd-%dh-%dm", t.Day(), t.Hour(), t.Minute()),
	}
}

func GenerateNameSuffix() string {
	t := time.Now()
	return fmt.Sprintf("%dd-%dh-%dm", t.Day(), t.Hour(), t.Minute())
}
