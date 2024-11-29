package capabilities

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func HealthCheck(client *gophercloud.ServiceClient) (r HealthCheckResult) {
	resp, err := client.Get(client.ServiceURL(), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))

	return
}

func GetImageCapabilities(client *gophercloud.ServiceClient, imageID string) (r ImageCapabilitiesResult) {
	resp, err := client.Get(imageCapabilitiesURL(client, imageID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))

	return
}
