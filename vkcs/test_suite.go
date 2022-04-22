package vkcs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/stretchr/testify/mock"
)

const testAccURL = "https://acctest.mcs.ru"

// dummyConfig is mock for Config
type dummyConfig struct {
	mock.Mock
}

var _ configer = &dummyConfig{}

// LoadAndValidate ...
func (d *dummyConfig) LoadAndValidate() error {
	args := d.Called()
	return args.Error(0)
}

// IdentityV3Client is a mock client for identity requests.
func (d *dummyConfig) IdentityV3Client(region string) (ContainerClient, error) {
	args := d.Called(region)
	if r, ok := args.Get(0).(ContainerClient); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// ContainerInfraV1Client is a mock client for infra requests.
func (d *dummyConfig) ContainerInfraV1Client(region string) (ContainerClient, error) {
	args := d.Called(region)
	if r, ok := args.Get(0).(ContainerClient); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// DatabaseV1Client returns dummy DatabaseV1Client
func (d *dummyConfig) DatabaseV1Client(region string) (*gophercloud.ServiceClient, error) {
	// args := d.Called(region)
	// if r, ok := args.Get(0).(ContainerClient); ok {
	// 	return r.(*gophercloud.ServiceClient), args.Error(1)
	// }
	// return nil, args.Error(0)
	return nil, nil
}

func (d *dummyConfig) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d *dummyConfig) ComputeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d *dummyConfig) ImageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d *dummyConfig) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d *dummyConfig) NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d *dummyConfig) GetMutex() *mutexkv.MutexKV {
	return nil
}

// GetRegion is a dummy method to return region.
func (d *dummyConfig) GetRegion() string {
	args := d.Called()
	return args.String(0)
}

// ContainerClientFixture ...
type ContainerClientFixture struct {
	mock.Mock
}

// Get ...
func (c *ContainerClientFixture) Get(url string, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, jsonResponse, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		if err := json.NewDecoder(r.Body).Decode(jsonResponse); err != nil {
			return r, args.Error(1)
		}
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// Post ...
func (c *ContainerClientFixture) Post(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, jsonBody, jsonResponse, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		if err := json.NewDecoder(r.Body).Decode(jsonResponse); err != nil {
			return r, args.Error(1)
		}
		return r, args.Error(1)
	}
	return nil, args.Error(0)

}

// Patch ...
func (c *ContainerClientFixture) Patch(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, jsonBody, jsonResponse, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		if err := json.NewDecoder(r.Body).Decode(jsonResponse); err != nil {
			return r, args.Error(1)
		}
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// Delete ...
func (c *ContainerClientFixture) Delete(url string, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// Head ...
func (c *ContainerClientFixture) Head(url string, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// Put ...
func (c *ContainerClientFixture) Put(url string, jsonBody interface{}, jsonResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error) {
	args := c.Called(url, jsonBody, jsonResponse, opts)
	if r, ok := args.Get(0).(*http.Response); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(0)
}

// ServiceURL ...
func (c *ContainerClientFixture) ServiceURL(parts ...string) string {
	args := c.Called(parts)
	return args.String(0) + "/" + strings.Join(parts, "/")
}

// FakeBody is struct that implements ReadCloser interface; use it for http.Response.Body mock
type FakeBody struct {
	body   []byte
	length int
}

func newFakeBody(jsonBody map[string]interface{}) (*FakeBody, error) {
	marshaled, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}
	return &FakeBody{
		body:   marshaled,
		length: len(marshaled),
	}, nil
}

// Read ...
func (f *FakeBody) Read(p []byte) (n int, err error) {
	copy(p, f.body)
	return len(p), nil
}

// Close ...
func (f *FakeBody) Close() (err error) {
	return nil
}

func makeClusterCreateResponseFixture(uuid string) *http.Response {
	fakeBody, _ := newFakeBody(map[string]interface{}{"uuid": uuid})
	resp := &http.Response{
		Status:        "202 Accepted",
		StatusCode:    202,
		Body:          fakeBody,
		ContentLength: int64(fakeBody.length),
	}
	return resp
}

func makeClusterGetResponseFixture(clusterGetFixture map[string]interface{}, uuid string, s clusterStatus) *http.Response {
	newMap := map[string]interface{}{}
	for k, v := range clusterGetFixture {
		newMap[k] = v
	}
	newMap["uuid"] = uuid
	newMap["new_status"] = s
	fakeBody, _ := newFakeBody(newMap)
	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Body:          fakeBody,
		ContentLength: int64(fakeBody.length),
	}
	return resp
}

func makeClusterDeleteResponseFixture() *http.Response {
	return &http.Response{
		Status:     "202 Accepted",
		StatusCode: 202,
	}
}
