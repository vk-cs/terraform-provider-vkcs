package providerjson

import (
	"encoding/json"
	"os"
)

func DumpWithWrapper(wrapper *ProviderWrapper, data *ProviderJSON) error {
	if s, err := ProviderFromRaw(data); err != nil {
		return err
	} else {
		wrapper.ProviderSchema = s
	}

	if err := json.NewEncoder(os.Stdout).Encode(wrapper); err != nil {
		return err
	}

	return nil
}

func WriteWithWrapper(base, wrapper *ProviderWrapper, data *ProviderJSON, filename, providerVersion string) error {
	if s, err := ProviderFromRaw(data); err != nil {
		return err
	} else {
		if base != nil {
			addNewSince(base.ProviderSchema, s, providerVersion)
		}
		wrapper.ProviderSchema = s
	}

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
