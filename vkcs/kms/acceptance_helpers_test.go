package kms_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func testAccKMSPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("OS_PROJECT_ID") == "" {
		t.Fatal("OS_PROJECT_ID must be set for KMS acceptance tests")
	}
	opts := clients.ConfigOpts{}
	if _, err := opts.LoadAndValidate(); err != nil {
		t.Fatalf("Invalid KMS acceptance credentials: %v", err)
	}
}

// Verify deletion through the API, including resources created by data source tests.
func testAccKMSCheckDestroy(state *terraform.State) error {
	opts := clients.ConfigOpts{}
	config, err := opts.LoadAndValidate()
	if err != nil {
		return err
	}
	client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		return err
	}
	for _, rs := range state.RootModule().Resources {
		switch rs.Type {
		case "vkcs_kms_key":
			_, err = keys.Get(client, rs.Primary.ID).Extract()
		case "vkcs_kms_secret":
			_, err = secrets.Get(client, rs.Primary.ID).Extract()
		default:
			continue
		}
		if err == nil {
			return fmt.Errorf("%s %q still exists", rs.Type, rs.Primary.ID)
		}
		if !errutil.IsNotFound(err) {
			return err
		}
	}
	return nil
}
