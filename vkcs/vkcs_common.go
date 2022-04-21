package vkcs

import (
	"net/http"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
)

// ContainerClient is interface to work with gophercloud requests
type ContainerClient interface {
	Get(url string, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Post(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Patch(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Delete(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Head(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Put(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	ServiceURL(parts ...string) string
}

func getDBRequestOpts(codes ...int) *gophercloud.RequestOpts {
	reqOpts := &gophercloud.RequestOpts{
		OkCodes: codes,
	}
	if len(codes) != 0 {
		reqOpts.OkCodes = codes
	}
	return reqOpts
}

// dateTimeWithoutTZFormat represents format of time used in dbaas
type dateTimeWithoutTZFormat struct {
	time.Time
}

// UnmarshalJSON is used to correctly unmarshal datetime fields
func (t *dateTimeWithoutTZFormat) UnmarshalJSON(b []byte) (err error) {
	layout := "2006-01-02T15:04:05"
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	t.Time, err = time.Parse(layout, s)
	return
}

type optsBuilder interface {
	Map() (map[string]interface{}, error)
}
