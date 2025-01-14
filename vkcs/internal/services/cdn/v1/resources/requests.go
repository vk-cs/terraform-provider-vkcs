package resources

import "github.com/gophercloud/gophercloud"

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts specifies attributes used to create a CDN resource.
type CreateOpts struct {
	Active             bool                   `json:"active"`
	CNAME              string                 `json:"cname"`
	Enabled            *bool                  `json:"enabled,omitempty"`
	LogTarget          string                 `json:"logTarget,omitempty"`
	Options            *ResourceOptions       `json:"options,omitempty"`
	Origin             string                 `json:"origin,omitempty"`
	OriginGroup        int                    `json:"originGroup,omitempty"`
	OriginProtocol     ResourceOriginProtocol `json:"originProtocol,omitempty"`
	SecondaryHostnames []string               `json:"secondaryHostnames"`
	SSLAutomated       bool                   `json:"ssl_automated,omitempty"`
	SSLData            int                    `json:"sslData,omitempty"`
	SSLEnabled         bool                   `json:"sslEnabled,omitempty"`
}

// Map builds a request body from a CreateOpts structure.
func (opts CreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements a resource create request.
func Create(client *gophercloud.ServiceClient, projectID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(resourcesURL(client, projectID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about a resource, given its ID.
func Get(client *gophercloud.ServiceClient, projectID string, id int) (r GetResult) {
	resp, err := client.Get(resourceURL(client, projectID, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetShielding returns information about origin shielding options
// applied to the resource.
func GetShielding(client *gophercloud.ServiceClient, projectID string, id int) (r GetShieldingResult) {
	resp, err := client.Get(shieldingURL(client, projectID, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateOpts specifies attributes used to update a CDN resource.
type UpdateOpts struct {
	Active             *bool                  `json:"active,omitempty"`
	Enabled            *bool                  `json:"enabled,omitempty"`
	LogTarget          string                 `json:"logTarget,omitempty"`
	Options            *ResourceOptions       `json:"options,omitempty"`
	Origin             string                 `json:"origin,omitempty"`
	OriginGroup        int                    `json:"originGroup,omitempty"`
	OriginProtocol     ResourceOriginProtocol `json:"originProtocol,omitempty"`
	SecondaryHostnames []string               `json:"secondaryHostnames"`
	SSLAutomated       *bool                  `json:"ssl_automated,omitempty"`
	SSLData            int                    `json:"sslData,omitempty"`
	SSLEnabled         *bool                  `json:"sslEnabled,omitempty"`
}

// Map builds a request body from a UpdateOpts structure.
func (opts UpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Update implements a resource update request.
func Update(client *gophercloud.ServiceClient, projectID string, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(resourceURL(client, projectID, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateShieldingOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateShieldingOpts specifies attributes used to update origin shielding settings.
type UpdateShieldingOpts struct {
	ShieldingPop *int `json:"shielding_pop"`
}

// Map builds a request body from a UpdateShieldingOpts structure.
func (opts UpdateShieldingOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// UpdateShielding implements a request to update origin shielding settings.
func UpdateShielding(client *gophercloud.ServiceClient, projectID string, resourceID int, opts UpdateShieldingOptsBuilder) (r UpdateShieldingResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(shieldingURL(client, projectID, resourceID), &b, nil, &gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// IssueLetsEncrypt implements a request to issue a Let's Encrypt certificate.
func IssueLetsEncrypt(client *gophercloud.ServiceClient, projectID string, resourceID int) (r IssueLetsEncryptResult) {
	resp, err := client.Post(issueLetsEncryptURL(client, projectID, resourceID), nil, nil, &gophercloud.RequestOpts{
		OkCodes:      []int{201},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetLetsEncryptStatus returns details on Let's Encrypt certificate issuance.
func GetLetsEncryptStatus(client *gophercloud.ServiceClient, projectID string, id int) (r GetLetsEncryptStatusResult) {
	resp, err := client.Get(getLetsEncryptStatusURL(client, projectID, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type PrefetchContentOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// PrefetchContentOpts specifies attributes used to prefetch content for a CDN resource.
type PrefetchContentOpts struct {
	Paths []string `json:"paths"`
}

// Map builds a request body from a PrefetchContentOpts structure.
func (opts PrefetchContentOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// PrefetchContent implements a request to prefetch content for a CDN resource.
func PrefetchContent(client *gophercloud.ServiceClient, projectID string, resourceID int, opts PrefetchContentOptsBuilder) (r PrefetchContentResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(prefetchContentURL(client, projectID, resourceID), &b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a resource delete request.
func Delete(client *gophercloud.ServiceClient, projectID string, id int) (r DeleteResult) {
	resp, err := client.Delete(resourceURL(client, projectID, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
