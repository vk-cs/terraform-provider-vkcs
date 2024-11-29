package capabilities

import (
	"github.com/gophercloud/gophercloud"
)

type HealthCheckResult struct {
	gophercloud.ErrResult
}

type ImageCapabilities struct {
	CapabilityVersions []struct {
		Capability struct {
			Description string `json:"description"`
			Name        string `json:"name"`
		} `json:"capability"`
		Default bool   `json:"default"`
		ID      string `json:"id"`
		Os      struct {
			Dist    string `json:"dist"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"os"`
		Package string `json:"package"`
		Version string `json:"version"`
	} `json:"capability_versions"`
}

// ImageCapabilitiesResult is the result of a get image capabilities request. Call its Extract method
// to interpret the result as ImageCapabilities
type ImageCapabilitiesResult struct {
	gophercloud.Result
}

// Extract extracts ImageCapabilities from a ImageCapabilitiesResult.
func (r ImageCapabilitiesResult) Extract() (ImageCapabilities, error) {
	var capabilities ImageCapabilities
	err := r.ExtractInto(&capabilities)
	return capabilities, err
}
