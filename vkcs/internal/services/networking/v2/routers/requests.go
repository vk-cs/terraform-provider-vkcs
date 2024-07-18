package routers

import (
	"errors"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
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

const creationRetriesNum = 3
const creationRetryDelay = 3 * time.Second

func needRetryOnNetworkingRouterCreationError(err error) bool {
	if errutil.IsNotFound(err) {
		return true
	}

	var http409Err gophercloud.ErrDefault409

	if errors.As(err, &http409Err) {
		neutronError, decodeErr := inetworking.DecodeNeutronError(http409Err.ErrUnexpectedResponseCode.Body)
		if decodeErr != nil {
			log.Printf("[DEBUG] failed to decode a neutron error: %s", decodeErr)
			return false
		}
		if neutronError.Type == "NeutronDbObjectDuplicateEntry" {
			time.Sleep(creationRetryDelay)
			return true
		}

		return false
	}

	return false
}

func CreateWithRetry(c *gophercloud.ServiceClient, opts routers.CreateOptsBuilder) routers.CreateResult {
	timer := time.NewTimer(time.Nanosecond)
	defer timer.Stop()
	for i := 0; i < creationRetriesNum; {
		select {
		case <-timer.C:
			r := routers.Create(c, opts)
			if r.Err == nil {
				return r
			}
			if !needRetryOnNetworkingRouterCreationError(r.Err) {
				r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
				return r
			}
			timer.Reset(creationRetryDelay)
			i++
		case <-c.Context.Done():
			res := routers.CreateResult{}
			res.Err = gophercloud.ErrTimeOut{}
			return res
		}
	}

	return routers.CreateResult{}
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
