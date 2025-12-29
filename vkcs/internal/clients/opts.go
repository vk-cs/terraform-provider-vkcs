package clients

import (
	"fmt"
	"os"

	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/version"
)

const (
	defaultEnvPrefix                     = "OS_"
	defaultIdentityEndpoint              = "https://infra.mail.ru/identity/v3/"
	defaulUserDomainName                 = "users"
	defaulRegionName                     = "RegionOne"
	defaultContainerInfraAPIMicroVersion = "1.33"
)

type ConfigOpts struct {
	EnvPrefix                    string
	Token                        string
	Username                     string
	Password                     string
	ProjectID                    string
	Region                       string
	IdentityEndpoint             string
	UserDomainID                 string
	UserDomainName               string
	EndpointType                 string
	EndpointOverrides            map[string]any
	TerraformVersion             string
	FrameworkVersion             string
	ContainerInfraV1MicroVersion string
	SkipAuth                     bool
}

// LoadAndValidate applies environment variables to the config, sets defaults, and validates
// provided authentication options.
func (o *ConfigOpts) LoadAndValidate() (Config, error) {
	if o.EnvPrefix == "" {
		o.EnvPrefix = defaultEnvPrefix
	}

	getEnv := func(key string) string {
		return os.Getenv(o.EnvPrefix + key)
	}

	if o.Username == "" {
		o.Username = getEnv("USERNAME")
	}

	if o.Password == "" {
		o.Password = getEnv("PASSWORD")
	}

	if o.ProjectID == "" {
		o.ProjectID = getEnv("PROJECT_ID")
	}

	if o.Region == "" {
		o.Region = getEnv("REGION_NAME")
	}

	if o.IdentityEndpoint == "" {
		o.IdentityEndpoint = getEnv("AUTH_URL")
	}

	if o.UserDomainID == "" {
		o.UserDomainID = getEnv("USER_DOMAIN_ID")
	}

	if o.UserDomainName == "" {
		o.UserDomainName = getEnv("USER_DOMAIN_NAME")
	}

	if o.EndpointType == "" {
		o.EndpointType = getEnv("INTERFACE")
	}

	if o.Token == "" {
		o.Token = getEnv("AUTH_TOKEN")
	}

	if o.TerraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		o.TerraformVersion = "0.11+compatible"
	}

	if o.ContainerInfraV1MicroVersion == "" {
		o.ContainerInfraV1MicroVersion = defaultContainerInfraAPIMicroVersion
	}

	if o.Region == "" {
		o.Region = defaulRegionName
	}

	if o.IdentityEndpoint == "" {
		o.IdentityEndpoint = defaultIdentityEndpoint
	}

	if o.UserDomainName == "" {
		o.UserDomainName = defaulUserDomainName
	}

	if o.UserDomainID != "" {
		o.UserDomainName = ""
	}

	authCfg := o.ToAuthConfig()
	if err := authCfg.LoadAndValidate(); err != nil {
		return nil, err
	}

	authCfg.OsClient.UserAgent.Prepend(fmt.Sprintf("VKCS Terraform Provider/%s", version.ProviderVersion))
	authCfg.OsClient.RetryFunc = retryFunc

	return &config{
		Config:                       authCfg,
		envPrefix:                    o.EnvPrefix,
		containerInfraV1MicroVersion: o.ContainerInfraV1MicroVersion,
		skipAuth:                     o.SkipAuth,
	}, nil
}

func (o *ConfigOpts) ToAuthConfig() auth.Config {
	cfg := auth.Config{
		IdentityEndpoint:  o.IdentityEndpoint,
		Username:          o.Username,
		Password:          o.Password,
		TenantID:          o.ProjectID,
		UserDomainID:      o.UserDomainID,
		UserDomainName:    o.UserDomainName,
		Region:            o.Region,
		Token:             o.Token,
		MaxRetries:        maxRetriesCount,
		MutexKV:           mutexkv.NewMutexKV(),
		EndpointOverrides: o.EndpointOverrides,
		SDKVersion:        o.FrameworkVersion,
		TerraformVersion:  o.TerraformVersion,
	}

	if o.Token == "" {
		cfg.AllowReauth = true
	} else {
		cfg.DelayedAuth = true
	}

	return cfg
}
