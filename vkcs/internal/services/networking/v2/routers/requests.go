package routers

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	creationRetriesNum = 3
	creationRetryDelay = 3 * time.Second
	creationTimeout    = 3 * time.Minute

	NeutronErrTypeDbObjectDuplicateEntry = "NeutronDbObjectDuplicateEntry"
)

// RouterCreateOpts represents the attributes used when creating a new router.
type RouterCreateOpts struct {
	routers.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToRouterCreateMap casts a CreateOpts struct to a map.
// It overrides routers.ToRouterCreateMap to add the ValueSpecs field.
func (opts RouterCreateOpts) ToRouterCreateMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "router")
}

type ExpectedNeutronError struct {
	ErrCode int
	ErrType string
}

func retryNeutronError(actual error, retryableErrors []ExpectedNeutronError) bool {
	for _, expectedErr := range retryableErrors {
		gophercloudErr, ok := errutil.As(actual, expectedErr.ErrCode)
		if !ok {
			continue
		}
		if expectedErr.ErrType == "" {
			return true
		}

		neutronError, decodeErr := inetworking.DecodeNeutronError(gophercloudErr.Body)
		if decodeErr != nil {
			continue
		}
		if expectedErr.ErrType == neutronError.Type {
			return true
		}
	}

	return false
}

func CreateWithRetry(c *gophercloud.ServiceClient, opts routers.CreateOptsBuilder) routers.CreateResult {
	retryableErrors := []ExpectedNeutronError{
		{ErrCode: 404},
		{ErrCode: 409, ErrType: NeutronErrTypeDbObjectDuplicateEntry},
	}

	var r routers.CreateResult
	var count int

	ctx := context.Background()
	if c.Context != nil {
		ctx = c.Context
	}

	createErr := retry.RetryContext(ctx, creationTimeout, func() *retry.RetryError {
		r = routers.Create(c, opts)
		if r.Err != nil {
			if count++; count >= creationRetriesNum {
				return retry.NonRetryableError(r.Err)
			}

			if retryNeutronError(r.Err, retryableErrors) {
				time.Sleep(creationRetryDelay)
				return retry.RetryableError(r.Err)
			}

			return retry.NonRetryableError(r.Err)
		}

		return nil
	})

	if createErr != nil {
		r.Err = util.ErrorWithRequestID(createErr, r.Header.Get(util.RequestIDHeader))
	}

	return r
}

func Get(c *gophercloud.ServiceClient, id string) routers.GetResult {
	r := routers.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts routers.UpdateOptsBuilder) routers.UpdateResult {
	r := routers.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) routers.DeleteResult {
	r := routers.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func AddInterface(c *gophercloud.ServiceClient, id string, opts routers.AddInterfaceOptsBuilder) routers.InterfaceResult {
	r := routers.AddInterface(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func RemoveInterface(c *gophercloud.ServiceClient, id string, opts routers.RemoveInterfaceOptsBuilder) routers.InterfaceResult {
	r := routers.RemoveInterface(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
