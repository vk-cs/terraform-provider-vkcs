package resources

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any resource result as a Resource, if possible.
func (r commonResult) Extract() (*Resource, error) {
	var res Resource
	err := r.ExtractInto(&res)
	return &res, err
}

// CreateResult is the result of a create request. Call its Extract method
// to interpret a result as a Resource.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as a Resource.
type GetResult struct {
	commonResult
}

type commonShieldingResult struct {
	gophercloud.Result
}

// Extract interprets any resource shielding result as a ResourceShielding,
// if possible.
func (r commonShieldingResult) Extract() (*ResourceShielding, error) {
	var res ResourceShielding
	err := r.ExtractInto(&res)
	return &res, err
}

// GetShieldingResult is the result of a get request. Call its Extract method
// to interpret a result as a ResourceShielding.
type GetShieldingResult struct {
	commonShieldingResult
}

// UpdateResult is the result of a delete request. Call its Extract method
// to interpret a result as a Resource.
type UpdateResult struct {
	commonResult
}

// UpdateShieldingResult is the result of a request to update origin shielding
// settings. Call its Extract method to interpret a result as a ResourceShielding.
type UpdateShieldingResult struct {
	commonShieldingResult
}

// IssueLetsEncryptResult is the result of a request to issue Let's Encrypt
// certificate. Call its ExtractErr method to determine if a request succeeded or failed.
type IssueLetsEncryptResult struct {
	gophercloud.ErrResult
}

// PrefetchContentResult is the result of a request to prefetch content for a CDN
// resource. Call its ExtractErr method to determine if a request succeeded or failed.
type PrefetchContentResult struct {
	gophercloud.ErrResult
}

// DeleteResult is the result of a delete request. Call its ExtractErr method
// to determine if a request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Resource represents a CDN resource.
type Resource struct {
	Active             bool                   `json:"active"`
	Client             int                    `json:"client"`
	CNAME              string                 `json:"cname"`
	CompanyName        string                 `json:"companyName"`
	Created            string                 `json:"created"`
	Deleted            bool                   `json:"deleted"`
	Enabled            bool                   `json:"enabled"`
	ID                 int                    `json:"id"`
	Options            ResourceOptions        `json:"options"`
	OriginGroup        int                    `json:"originGroup"`
	OriginGroupName    string                 `json:"originGroup_name"`
	OriginProtocol     ResourceOriginProtocol `json:"originProtocol"`
	PresetApplied      bool                   `json:"preset_applied"`
	ProxySSLEnabled    bool                   `json:"proxy_ssl_enabled"`
	SecondaryHostnames []string               `json:"secondaryHostnames"`
	Shielded           bool                   `json:"shielded"`
	SSLData            int                    `json:"sslData"`
	SSLEnabled         bool                   `json:"sslEnabled"`
	SSLAutomated       bool                   `json:"ssl_automated"`
	Status             string                 `json:"status"`
	Updated            string                 `json:"updated"`
	VPEnabled          bool                   `json:"vp_enabled"`
}

// ResourceShielding represents origin shielding options applied to the resource.
type ResourceShielding struct {
	ShieldingPop *int `json:"shielding_pop"`
}
