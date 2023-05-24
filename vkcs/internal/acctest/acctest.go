package acctest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
)

var (
	OsFlavorName       = os.Getenv("OS_FLAVOR_NAME")
	OsNewFlavorName    = os.Getenv("OS_NEW_FLAVOR_NAME")
	OsImageName        = os.Getenv("OS_IMAGE_NAME")
	OsRegionName       = os.Getenv("OS_REGION_NAME")
	OsProjectID        = os.Getenv("OS_PROJECT_ID")
	OsExtNetName       = os.Getenv("OS_EXT_NET_NAME")
	OsAvailabilityZone = os.Getenv("OS_AVAILABILITY_ZONE")
	OsVolumeType       = os.Getenv("OS_VOLUME_TYPE")
	// Kubernetes-related environment variables
	OsFlavorID        = os.Getenv("OS_FLAVOR_ID")
	OsNewFlavorID     = os.Getenv("OS_NEW_FLAVOR_ID")
	OsNetworkID       = os.Getenv("OS_NETWORK_ID")
	OsSubnetworkID    = os.Getenv("OS_SUBNETWORK_ID")
	OsKeypairName     = os.Getenv("OS_KEYPAIR_NAME")
	ClusterTemplateID = os.Getenv("CLUSTER_TEMPLATE_ID")
)

var AccTestValues map[string]string = map[string]string{
	"BaseNetwork":      AccTestBaseNetwork,
	"BaseExtNetwork":   AccTestBaseExtNetwork(),
	"BaseImage":        AccTestBaseImage(),
	"BaseFlavor":       AccTestBaseFlavor(),
	"BaseNewFlavor":    AccTestBaseNewFlavor(),
	"AvailabilityZone": OsAvailabilityZone,
	"VolumeType":       OsVolumeType,
	"FlavorName":       OsFlavorName,
	"NewFlavorName":    OsNewFlavorName,
	"ImageName":        OsImageName,
	"ExtNetName":       OsExtNetName,
	"ProjectID":        OsProjectID,
}

var AccTestProviders map[string]func() (*schema.Provider, error)
var AccTestProvider *schema.Provider
var AccTestProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

func init() {
	AccTestProvider = provider.Provider()
	AccTestProviders = map[string]func() (*schema.Provider, error){
		"vkcs": func() (*schema.Provider, error) {
			return AccTestProvider, nil
		},
	}
	AccTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"vkcs": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(provider.New()),
				func() tfprotov6.ProviderServer {
					server, _ := tf5to6server.UpgradeServer(
						ctx,
						provider.Provider().GRPCProvider,
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
		"OS_VOLUME_TYPE":       OsVolumeType,
		"OS_AVAILABILITY_ZONE": OsAvailabilityZone,
		"OS_FLAVOR_NAME":       OsFlavorName,
		"OS_NEW_FLAVOR_NAME":   OsNewFlavorName,
		"OS_IMAGE_NAME":        OsImageName,
		"OS_EXT_NET_NAME":      OsExtNetName,
		"OS_PROJECT_ID":        OsProjectID,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
}

func AccTestPreCheckKubernetes(t *testing.T) {
	vars := map[string]interface{}{
		"CLUSTER_TEMPLATE_ID": ClusterTemplateID,
		"OS_FLAVOR_ID":        OsFlavorID,
		"OS_NETWORK_ID":       OsNetworkID,
		"OS_SUBNETWORK_ID":    OsSubnetworkID,
		"OS_KEYPAIR_NAME":     OsKeypairName,
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

func AccTestBaseImage() string {
	return fmt.Sprintf(`
	data "vkcs_images_image" "base" {
		name = "%s"
	}
`, OsImageName)
}

const AccTestBaseNetwork string = `

data "vkcs_networking_network" "extnet" {
	name = "ext-net"
  }
  
  resource "vkcs_networking_network" "base" {
	name           = "base-net"
	admin_state_up = true
  }
  
  resource "vkcs_networking_subnet" "base" {
	name       = "subnet_1"
	network_id = vkcs_networking_network.base.id
	cidr       = "192.168.199.0/24"
  }
  
  resource "vkcs_networking_router" "base" {
	name                = "base-router"
	admin_state_up      = true
	external_network_id = data.vkcs_networking_network.extnet.id
  }
  
  resource "vkcs_networking_router_interface" "base" {
	router_id = vkcs_networking_router.base.id
	subnet_id = vkcs_networking_subnet.base.id
  }
`

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
