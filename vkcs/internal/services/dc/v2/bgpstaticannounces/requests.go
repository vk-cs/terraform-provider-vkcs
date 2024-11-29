package bgpstaticannounces

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type BGPStaticAnnounceCreate struct {
	BGPStaticAnnounce *CreateOpts `json:"dc_bgp_static_announce"`
}

type CreateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	DCBGPID     string `json:"dc_bgp_id"`
	Network     string `json:"network"`
	Gateway     string `json:"gateway"`
	Enabled     *bool  `json:"enabled,omitempty"`
}

func (opts *BGPStaticAnnounceCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(bgpStaticAnnouncesURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(bgpStaticAnnounceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type BGPStaticAnnounceUpdate struct {
	BGPStaticAnnounce *UpdateOpts `json:"dc_bgp_static_announce"`
}

type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
}

func (opts *BGPStaticAnnounceUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(bgpStaticAnnounceURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(bgpStaticAnnounceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
