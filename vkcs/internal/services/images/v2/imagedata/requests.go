package imagedata

import (
	"io"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Upload(client *gophercloud.ServiceClient, id string, data io.Reader) imagedata.UploadResult {
	r := imagedata.Upload(client, id, data)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
