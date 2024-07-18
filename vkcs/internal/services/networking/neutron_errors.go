package networking

import (
	"encoding/json"
)

type NeutronError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
}

func DecodeNeutronError(body []byte) (*NeutronError, error) {
	neutronErr := &struct {
		NeutronError NeutronError
	}{}
	if err := json.Unmarshal(body, neutronErr); err != nil {
		return nil, err
	}

	return &neutronErr.NeutronError, nil
}
