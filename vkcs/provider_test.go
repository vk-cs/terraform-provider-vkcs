package vkcs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/pathorcontents"
)

var (
	osFlavorName       = os.Getenv("OS_FLAVOR_NAME")
	osNewFlavorName    = os.Getenv("OS_NEW_FLAVOR_NAME")
	osImageName        = os.Getenv("OS_IMAGE_NAME")
	osRegionName       = os.Getenv("OS_REGION_NAME")
	osProjectID        = os.Getenv("OS_PROJECT_ID")
	osExtNetName       = os.Getenv("OS_EXT_NET_NAME")
	osAvailabilityZone = os.Getenv("OS_AVAILABILITY_ZONE")
	osVolumeType       = os.Getenv("OS_VOLUME_TYPE")
	// Kubernetes-related environment variables
	osFlavorID        = os.Getenv("OS_FLAVOR_ID")
	osNewFlavorID     = os.Getenv("OS_NEW_FLAVOR_ID")
	osNetworkID       = os.Getenv("OS_NETWORK_ID")
	osSubnetworkID    = os.Getenv("OS_SUBNETWORK_ID")
	osKeypairName     = os.Getenv("OS_KEYPAIR_NAME")
	clusterTemplateID = os.Getenv("CLUSTER_TEMPLATE_ID")
)

var testAccValues map[string]string = map[string]string{
	"BaseNetwork":      testAccBaseNetwork,
	"BaseExtNetwork":   testAccBaseExtNetwork(),
	"BaseImage":        testAccBaseImage(),
	"BaseFlavor":       testAccBaseFlavor(),
	"BaseNewFlavor":    testAccBaseNewFlavor(),
	"AvailabilityZone": osAvailabilityZone,
	"VolumeType":       osVolumeType,
	"FlavorName":       osFlavorName,
	"NewFlavorName":    osNewFlavorName,
	"ImageName":        osImageName,
	"ExtNetName":       osExtNetName,
	"ProjectID":        osProjectID,
}

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"vkcs": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	vars := map[string]interface{}{
		"OS_VOLUME_TYPE":       osVolumeType,
		"OS_AVAILABILITY_ZONE": osAvailabilityZone,
		"OS_FLAVOR_NAME":       osFlavorName,
		"OS_NEW_FLAVOR_NAME":   osNewFlavorName,
		"OS_IMAGE_NAME":        osImageName,
		"OS_EXT_NET_NAME":      osExtNetName,
		"OS_PROJECT_ID":        osProjectID,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
}

func testAccPreCheckKubernetes(t *testing.T) {
	vars := map[string]interface{}{
		"CLUSTER_TEMPLATE_ID": clusterTemplateID,
		"OS_FLAVOR_ID":        osFlavorID,
		"OS_NETWORK_ID":       osNetworkID,
		"OS_SUBNETWORK_ID":    osSubnetworkID,
		"OS_KEYPAIR_NAME":     osKeypairName,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
}

func testAccAuthFromEnv() (configer, error) {
	tenantID := os.Getenv("OS_TENANT_ID")
	if tenantID == "" {
		tenantID = os.Getenv("OS_PROJECT_ID")
	}

	tenantName := os.Getenv("OS_TENANT_NAME")
	if tenantName == "" {
		tenantName = os.Getenv("OS_PROJECT_NAME")
	}

	config := &config{
		auth.Config{
			CACertFile:        os.Getenv("OS_CACERT"),
			ClientCertFile:    os.Getenv("OS_CERT"),
			ClientKeyFile:     os.Getenv("OS_KEY"),
			Cloud:             os.Getenv("OS_CLOUD"),
			DefaultDomain:     os.Getenv("OS_DEFAULT_DOMAIN"),
			DomainID:          os.Getenv("OS_DOMAIN_ID"),
			DomainName:        os.Getenv("OS_DOMAIN_NAME"),
			EndpointType:      os.Getenv("OS_ENDPOINT_TYPE"),
			IdentityEndpoint:  os.Getenv("OS_AUTH_URL"),
			Password:          os.Getenv("OS_PASSWORD"),
			ProjectDomainID:   os.Getenv("OS_PROJECT_DOMAIN_ID"),
			ProjectDomainName: os.Getenv("OS_PROJECT_DOMAIN_NAME"),
			Region:            os.Getenv("OS_REGION_NAME"),
			Token:             os.Getenv("OS_TOKEN"),
			TenantID:          tenantID,
			TenantName:        tenantName,
			UserDomainID:      os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName:    os.Getenv("OS_USER_DOMAIN_NAME"),
			Username:          os.Getenv("OS_USERNAME"),
			UserID:            os.Getenv("OS_USER_ID"),
			MutexKV:           mutexkv.NewMutexKV(),
		},
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return config, nil
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

// Steps for configuring OpenStack with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping VKCS SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping VKCS CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}

	diag := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("unexpected err when specifying VKCS CA by file: %v", diag)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping VKCS SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping VKCS CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}

	diag := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("Unexpected err when specifying VKCS CA by string: %v", diag)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping VKCS SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping VKCS client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}

	diag := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("unexpected err when specifying VKCS Client keypair by file: %v", diag)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping VKCS SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping VKCS client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}

	diag := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("unexpected err when specifying VKCS Client keypair by contents: %v", diag)
	}
}

func envVarContents(varName string) (string, error) {
	// TODO(irlndts): the function is deprecated, replace it.
	// nolint:staticcheck
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", varName)
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}

func testAccBaseExtNetwork() string {
	return fmt.Sprintf(`
	data "vkcs_networking_network" "extnet" {
		name = "%s"
	  }
	`, osExtNetName)
}

func testAccBaseFlavor() string {
	return fmt.Sprintf(`
	data "vkcs_compute_flavor" "base" {
		name = "%s"
	}
`, osFlavorName)
}

func testAccBaseNewFlavor() string {
	return fmt.Sprintf(`
	data "vkcs_compute_flavor" "base" {
		name = "%s"
	}
`, osNewFlavorName)
}

func testAccBaseImage() string {
	return fmt.Sprintf(`
	data "vkcs_images_image" "base" {
		name = "%s"
	}
`, osImageName)
}

const testAccBaseNetwork string = `

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
	ip_version = 4
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

func testAccRenderConfig(testConfig string, values ...map[string]string) string {
	t := template.Must(template.New("acc").Option("missingkey=error").Parse(testConfig))
	buf := &bytes.Buffer{}

	tmplValues := map[string]string{}
	copyToMap(&tmplValues, &testAccValues)
	if len(values) > 0 {
		copyToMap(&tmplValues, &values[0])
	}

	_ = t.Execute(buf, tmplValues)

	return buf.String()
}
