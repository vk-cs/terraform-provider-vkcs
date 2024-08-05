package routers

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	routerCreationRetriesNum = 3
	routerCreationRetryDelay = 3 * time.Second
	routerCreationTimeout    = 3 * time.Minute
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

func Create(c *gophercloud.ServiceClient, opts routers.CreateOptsBuilder) routers.CreateResult {
	retryableErrors := []inetworking.ExpectedNeutronError{
		{
			ErrCode: 404,
		},
		{
			ErrCode: 409,
			ErrType: inetworking.NeutronErrTypeDBObjectDuplicateEntry,
		},
	}

	ctx := context.Background()
	if c.Context != nil {
		ctx = c.Context
	}

	var r routers.CreateResult
	createFunc := func() error {
		r = routers.Create(c, opts)
		return r.Err
	}

	createErr := inetworking.CreateResourceWithRetry(ctx, createFunc, retryableErrors,
		routerCreationRetriesNum, routerCreationRetryDelay, routerCreationTimeout)

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
