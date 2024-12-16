package securitypolicytemplates

import (
	"encoding/base64"

	"github.com/gophercloud/gophercloud"
)

type securityPolicyTemplatesResult struct {
	gophercloud.Result
}

type SecurityPolicyTemplate struct {
	UUID                string `json:"uuid"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	SettingsDescription string `json:"settings_description"`
	Version             string `json:"version"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type SecurityPolicyTemplates struct {
	SecurityPolicyTemplates []SecurityPolicyTemplate `json:"security_policy"`
}

// Extract parses result into params for security policy templates.
func (r securityPolicyTemplatesResult) Extract() ([]SecurityPolicyTemplate, error) {
	var s *SecurityPolicyTemplates
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	for i, t := range s.SecurityPolicyTemplates {
		if t.SettingsDescription != "" {
			SettingsDescriptionDecoded, err := base64.StdEncoding.DecodeString(t.SettingsDescription)
			if err != nil {
				return nil, err
			}
			t.SettingsDescription = string(SettingsDescriptionDecoded)
			s.SecurityPolicyTemplates[i] = t
		}
	}

	return s.SecurityPolicyTemplates, err
}
