package templater

import (
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func ListUsers(client *gophercloud.ServiceClient, projectID string) (r ListUsersResult) {
	resp, err := client.Get(fmt.Sprintf("%s/list", serviceUserURL(client, projectID)), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))

	return
}

// CreateUserOpts specifies the attributes to get a monitoring settings
type CreateUserOpts struct {
	ImageID      string   `json:"image_id"`
	Capabilities []string `json:"capabilities"`
}

// Map formats a CreateUserOpts structure into a request body
func (opts CreateUserOpts) Map() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

func CreateUser(client *gophercloud.ServiceClient, projectID string, opts CreateUserOpts) (r CreateUserResult) {
	body, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(serviceUserURL(client, projectID), &body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusCreated},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))

	return
}

func DeleteUser(client *gophercloud.ServiceClient, projectID string, userID string) (r gophercloud.ErrResult) {
	resp, err := client.Delete(fmt.Sprintf("%s/%s", serviceUserURL(client, projectID), userID), &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusAccepted},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))

	return
}
