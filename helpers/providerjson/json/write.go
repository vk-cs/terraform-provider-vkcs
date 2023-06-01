package json

import (
	"encoding/json"
	"os"

	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func DumpWithWrapper(wrapper *schema.ProviderWrapper) error {
	if err := json.NewEncoder(os.Stdout).Encode(wrapper); err != nil {
		return err
	}

	return nil
}

func WriteWithWrapper(wrapper *schema.ProviderWrapper, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(wrapper); err != nil {
		return err
	}

	return nil
}
