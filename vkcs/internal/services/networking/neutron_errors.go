package networking

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	NeutronErrTypeDBObjectDuplicateEntry = "NeutronDbObjectDuplicateEntry"
	NeutronErrIPAddressGenerationFailure = "IpAddressGenerationFailure"
	NeutronErrExternalIPAddressExhausted = "ExternalIpAddressExhausted"
)

type NeutronError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
}

func DecodeNeutronError(body []byte) (*NeutronError, error) {
	neutronErr := &struct {
		NeutronError NeutronError
	}{}
	if err := json.Unmarshal(body, neutronErr); err != nil {
		return nil, err
	}

	return &neutronErr.NeutronError, nil
}

type ExpectedNeutronError struct {
	ErrCode int
	ErrType string
}

func RetryNeutronError(actual error, retryableErrors []ExpectedNeutronError, retryUnknown bool) bool {
	for _, expectedErr := range retryableErrors {
		gophercloudErr, ok := errutil.As(actual, expectedErr.ErrCode)
		if !ok {
			continue
		}
		if expectedErr.ErrType == "" {
			return true
		}

		neutronError, decodeErr := DecodeNeutronError(gophercloudErr.Body)
		if decodeErr != nil {
			log.Printf("[DEBUG] failed to decode a neutron error: %s", decodeErr)
			if retryUnknown {
				return true
			}
			continue
		}

		if expectedErr.ErrType == neutronError.Type {
			return true
		}
	}

	return false
}

func CreateResourceWithRetry(ctx context.Context, createFunc func() error, retryableErrors []ExpectedNeutronError,
	creationRetriesNum int, creationRetryDelay time.Duration, creationTimeout time.Duration) error {
	count := 0
	createErr := retry.RetryContext(ctx, creationTimeout, func() *retry.RetryError {
		err := createFunc()
		if err != nil {
			if count++; count >= creationRetriesNum {
				return retry.NonRetryableError(err)
			}

			if RetryNeutronError(err, retryableErrors, false) {
				time.Sleep(creationRetryDelay)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})

	return createErr
}
