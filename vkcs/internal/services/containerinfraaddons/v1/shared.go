package v1

import (
	"encoding/base64"
	"encoding/json"
)

type Addon struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ChartVersion            string `json:"chart_version"`
	ChartName               string `json:"chart_name"`
	ValuesTemplate          string `json:"values_template"`
	Description             string `json:"description"`
	CreatedAt               string `json:"created_at"`
	ShortDescription        string `json:"short_description"`
	InstallationInstruction string `json:"installation_instruction"`
	Purpose                 string `json:"purpose"`
	Changelog               string `json:"changelog"`
}

func (a *Addon) UnmarshalJSON(b []byte) error {
	type tmp Addon
	var s tmp

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	values := make([]byte, base64.StdEncoding.DecodedLen(len(s.ValuesTemplate)))
	n, err := base64.StdEncoding.Decode(values, []byte(s.ValuesTemplate))
	if err != nil {
		return err
	}

	s.ValuesTemplate = string(values[:n])
	*a = Addon(s)

	return nil
}

type Payload struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
