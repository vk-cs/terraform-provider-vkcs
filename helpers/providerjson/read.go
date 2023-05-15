package providerjson

import (
	"encoding/json"
	"os"
)

func ReadWithWrapper(filename string) (*ProviderWrapper, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pw ProviderWrapper
	err = json.NewDecoder(f).Decode(&pw)
	if err != nil {
		return nil, err
	}
	return &pw, nil
}
