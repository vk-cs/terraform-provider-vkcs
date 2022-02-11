package vkcs

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/pathorcontents"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

var (
	osFlavorID = os.Getenv("OS_FLAVOR_ID")
	// osFlavorName = os.Getenv("OS_FLAVOR_NAME")
	osImageID    = os.Getenv("OS_IMAGE_ID")
	osNetworkID  = os.Getenv("OS_NETWORK_ID")
	osRegionName = os.Getenv("OS_REGION_NAME")
	osProjectID  = os.Getenv("OS_PROJECT_ID")
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"vkcs": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
		// "terraform-provider-openstack/openstack": func() (*schema.Provider, error) {
		// 	return testAccProvider, nil
		// },
	}
}

func testAccPreCheckCompute(t *testing.T) {
	vars := map[string]interface{}{
		"OS_FLAVOR_ID": osFlavorID,
		"OS_IMAGE_ID":  osImageID,
		// "OS_FLAVOR_NAME": osFlavorName,
		"OS_NETWORK_ID":  osNetworkID,
		"OS_REGION_NAME": osRegionName,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
}

// func TestProvider(t *testing.T) {
// 	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
// }

// func TestProvider_impl(t *testing.T) {
// 	var _ = Provider()
// }

// // Steps for configuring OpenStack with SSL validation are here:
// // https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
// func TestAccProvider_caCertFile(t *testing.T) {
// 	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
// 		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
// 	}
// 	if os.Getenv("OS_CACERT") == "" {
// 		t.Skip("OS_CACERT is not set; skipping OpenStack CA test.")
// 	}

// 	p := Provider()

// 	caFile, err := envVarFile("OS_CACERT")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(caFile)

// 	raw := map[string]interface{}{
// 		"cacert_file": caFile,
// 	}

// 	err = p.Configure(terraform.NewResourceConfigRaw(raw))
// 	if err != nil {
// 		t.Fatalf("unexpected err when specifying OpenStack CA by file: %s", err)
// 	}
// }

// func TestAccProvider_caCertString(t *testing.T) {
// 	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
// 		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
// 	}
// 	if os.Getenv("OS_CACERT") == "" {
// 		t.Skip("OS_CACERT is not set; skipping OpenStack CA test.")
// 	}

// 	p := Provider()

// 	caContents, err := envVarContents("OS_CACERT")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	raw := map[string]interface{}{
// 		"cacert_file": caContents,
// 	}

// 	err = p.Configure(terraform.NewResourceConfigRaw(raw))
// 	if err != nil {
// 		t.Fatalf("Unexpected err when specifying OpenStack CA by string: %s", err)
// 	}
// }

// func TestAccProvider_clientCertFile(t *testing.T) {
// 	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
// 		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
// 	}
// 	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
// 		t.Skip("OS_CERT or OS_KEY is not set; skipping OpenStack client SSL auth test.")
// 	}

// 	p := Provider()

// 	certFile, err := envVarFile("OS_CERT")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(certFile)
// 	keyFile, err := envVarFile("OS_KEY")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(keyFile)

// 	raw := map[string]interface{}{
// 		"cert": certFile,
// 		"key":  keyFile,
// 	}

// 	err = p.Configure(terraform.NewResourceConfigRaw(raw))
// 	if err != nil {
// 		t.Fatalf("unexpected err when specifying OpenStack Client keypair by file: %s", err)
// 	}
// }

// func TestAccProvider_clientCertString(t *testing.T) {
// 	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
// 		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
// 	}
// 	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
// 		t.Skip("OS_CERT or OS_KEY is not set; skipping OpenStack client SSL auth test.")
// 	}

// 	p := Provider()

// 	certContents, err := envVarContents("OS_CERT")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	keyContents, err := envVarContents("OS_KEY")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	raw := map[string]interface{}{
// 		"cert": certContents,
// 		"key":  keyContents,
// 	}

// 	err = p.Configure(terraform.NewResourceConfigRaw(raw))
// 	if err != nil {
// 		t.Fatalf("unexpected err when specifying OpenStack Client keypair by contents: %s", err)
// 	}
// }

// func envVarContents(varName string) (string, error) {
// 	// TODO(irlndts): the function is deprecated, replace it.
// 	// nolint:staticcheck
// 	contents, _, err := pathorcontents.Read(os.Getenv(varName))
// 	if err != nil {
// 		return "", fmt.Errorf("error reading %s: %s", varName, err)
// 	}
// 	return contents, nil
// }

// func envVarFile(varName string) (string, error) {
// 	contents, err := envVarContents(varName)
// 	if err != nil {
// 		return "", err
// 	}

// 	tmpFile, err := ioutil.TempFile("", varName)
// 	if err != nil {
// 		return "", fmt.Errorf("error creating temp file: %s", err)
// 	}
// 	if _, err := tmpFile.Write([]byte(contents)); err != nil {
// 		_ = os.Remove(tmpFile.Name())
// 		return "", fmt.Errorf("error writing temp file: %s", err)
// 	}
// 	if err := tmpFile.Close(); err != nil {
// 		_ = os.Remove(tmpFile.Name())
// 		return "", fmt.Errorf("error closing temp file: %s", err)
// 	}
// 	return tmpFile.Name(), nil
// }
