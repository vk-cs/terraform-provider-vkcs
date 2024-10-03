package resources

// ResourceStatus CDN resource status
type ResourceStatus string

const (
	ResourceStatusActive    string = "active"
	ResourceStatusProcessed string = "processed"
	ResourceStatusSuspended string = "suspended"
)

// ResourceOriginProtocol origin interaction protocol
type ResourceOriginProtocol string

const (
	ResourceOriginProtocolHTTP  ResourceOriginProtocol = "HTTP"
	ResourceOriginProtocolHTTPS ResourceOriginProtocol = "HTTPS"
	ResourceOriginProtocolMATCH ResourceOriginProtocol = "MATCH"
)

// ResourceOriginProtocolValues returns list of all possible values for
// a ResourceOriginProtocol enumeration.
func ResourceOriginProtocolValues() []string {
	return []string{
		string(ResourceOriginProtocolHTTP),
		string(ResourceOriginProtocolHTTPS),
		string(ResourceOriginProtocolMATCH),
	}
}

// ResourceAllowedHttpMethod allowed HTTP method
type ResourceAllowedHttpMethod string

const (
	ResourceAllowedHttpMethodGet     ResourceAllowedHttpMethod = "GET"
	ResourceAllowedHttpMethodPost    ResourceAllowedHttpMethod = "POST"
	ResourceAllowedHttpMethodHead    ResourceAllowedHttpMethod = "HEAD"
	ResourceAllowedHttpMethodOptions ResourceAllowedHttpMethod = "OPTIONS"
	ResourceAllowedHttpMethodPut     ResourceAllowedHttpMethod = "PUT"
	ResourceAllowedHttpMethodPatch   ResourceAllowedHttpMethod = "PATCH"
	ResourceAllowedHttpMethodDelete  ResourceAllowedHttpMethod = "DELETE"
)

// ResourceAllowedHttpMethodValues returns list of all possible values for
// a ResourceAllowedHttpMethod enumeration.
func ResourceAllowedHttpMethodValues() []string {
	return []string{
		string(ResourceAllowedHttpMethodGet),
		string(ResourceAllowedHttpMethodPost),
		string(ResourceAllowedHttpMethodHead),
		string(ResourceAllowedHttpMethodOptions),
		string(ResourceAllowedHttpMethodPut),
		string(ResourceAllowedHttpMethodPatch),
		string(ResourceAllowedHttpMethodDelete),
	}
}

type ResourceACLPolicyType string

const (
	ResourceACLPolicyTypeAllow ResourceACLPolicyType = "allow"
	ResourceACLPolicyTypeDeny  ResourceACLPolicyType = "deny"
)

// ResourceACLPolicyTypeValues returns list of all possible values for
// a ResourceACLPolicyType enumeration.
func ResourceACLPolicyTypeValues() []string {
	return []string{
		string(ResourceACLPolicyTypeAllow),
		string(ResourceACLPolicyTypeDeny),
	}
}

type ResourceOptionsBoolOption struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type ResourceOptionsStringOption struct {
	Enabled bool   `json:"enabled"`
	Value   string `json:"value"`
}

type ResourceOptionsStringListOption struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type ResourceOptionsStringMapOption struct {
	Enabled bool              `json:"enabled"`
	Value   map[string]string `json:"value"`
}

type ResourceOptionsAllowedHttpMethodsOption struct {
	Enabled bool                        `json:"enabled"`
	Value   []ResourceAllowedHttpMethod `json:"value"`
}

type ResourceOptionsEdgeCacheSettingsOption struct {
	CustomValues *map[string]string `json:"custom_values,omitempty"`
	Default      string             `json:"default,omitempty"`
	Enabled      bool               `json:"enabled"`
	Value        string             `json:"value,omitempty"`
}

type ResourceOptionsForceReturnOption struct {
	Body    string `json:"body,omitempty"`
	Code    int    `json:"code,omitempty"`
	Enabled bool   `json:"enabled"`
}

type ResourceOptionsACLOption struct {
	Enabled        bool                  `json:"enabled"`
	ExceptedValues []string              `json:"excepted_values"`
	PolicyType     ResourceACLPolicyType `json:"policy_type"`
}

type ResourceOptions struct {
	AllowedHttpMethods   *ResourceOptionsAllowedHttpMethodsOption `json:"allowedHttpMethods,omitempty"`
	BrotliCompression    *ResourceOptionsStringListOption         `json:"brotli_compression,omitempty"`
	BrowserCacheSettings *ResourceOptionsStringOption             `json:"browser_cache_settings,omitempty"`
	CORS                 *ResourceOptionsStringListOption         `json:"cors,omitempty"`
	EdgeCacheSettings    *ResourceOptionsEdgeCacheSettingsOption  `json:"edge_cache_settings,omitempty"`
	FetchCompressed      *ResourceOptionsBoolOption               `json:"fetch_compressed,omitempty"`
	ForceReturn          *ResourceOptionsForceReturnOption        `json:"force_return,omitempty"`
	ForwardHostHeader    *ResourceOptionsBoolOption               `json:"forward_host_header,omitempty"`
	GzipOn               *ResourceOptionsBoolOption               `json:"gzipOn,omitempty"`
	HostHeader           *ResourceOptionsStringOption             `json:"hostHeader,omitempty"`
	IgnoreQueryString    *ResourceOptionsBoolOption               `json:"ignoreQueryString,omitempty"`
	IgnoreCookie         *ResourceOptionsBoolOption               `json:"ignore_cookie,omitempty"`
	QueryParamsBlacklist *ResourceOptionsStringListOption         `json:"query_params_blacklist,omitempty"`
	QueryParamsWhitelist *ResourceOptionsStringListOption         `json:"query_params_whitelist,omitempty"`
	CountryACL           *ResourceOptionsACLOption                `json:"country_acl,omitempty"`
	ReferrerACL          *ResourceOptionsACLOption                `json:"referrer_acl,omitempty"`
	IpAddressACL         *ResourceOptionsACLOption                `json:"ip_address_acl,omitempty"`
	UserAgentACL         *ResourceOptionsACLOption                `json:"user_agent_acl,omitempty"`
	Slice                *ResourceOptionsBoolOption               `json:"slice,omitempty"`
	Stale                *ResourceOptionsStringListOption         `json:"stale,omitempty"`
	StaticHeaders        *ResourceOptionsStringMapOption          `json:"staticHeaders,omitempty"`
	StaticRequestHeaders *ResourceOptionsStringMapOption          `json:"staticRequestHeaders,omitempty"`
	Websockets           *ResourceOptionsBoolOption               `json:"websockets,omitempty"`
}
