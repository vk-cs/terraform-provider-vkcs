package json

import (
	"encoding/json"
	"os"

	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func ReadWithWrapper(filename string) (*schema.ProviderWrapper, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pw schema.ProviderWrapper
	err = json.NewDecoder(f).Decode(&pw)
	if err != nil {
		return nil, err
	}
	return &pw, nil
}
