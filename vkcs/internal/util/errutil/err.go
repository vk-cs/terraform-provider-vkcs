package errutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
)

const RequestIDHeader = "X-Openstack-Request-Id"

func Is(err error, errorCode int) bool {
	if err == nil {
		return false
	}

	switch errorCode {
	case 400:
		var httpErr gophercloud.ErrDefault400
		return errors.As(err, &httpErr)
	case 401:
		var httpErr gophercloud.ErrDefault401
		return errors.As(err, &httpErr)
	case 403:
		var httpErr gophercloud.ErrDefault403
		return errors.As(err, &httpErr)
	case 404:
		var httpErr gophercloud.ErrDefault404
		return errors.As(err, &httpErr)
	case 405:
		var httpErr gophercloud.ErrDefault405
		return errors.As(err, &httpErr)
	case 408:
		var httpErr gophercloud.ErrDefault408
		return errors.As(err, &httpErr)
	case 409:
		var httpErr gophercloud.ErrDefault409
		return errors.As(err, &httpErr)
	case 429:
		var httpErr gophercloud.ErrDefault429
		return errors.As(err, &httpErr)
	case 500:
		var httpErr gophercloud.ErrDefault500
		return errors.As(err, &httpErr)
	case 503:
		var httpErr gophercloud.ErrDefault503
		return errors.As(err, &httpErr)
	}

	var unknownErr gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &unknownErr) {
		return unknownErr.Actual == errorCode
	}
	return false
}

func As(err error, errorCode int) (*gophercloud.ErrUnexpectedResponseCode, bool) {
	if err == nil {
		return nil, false
	}

	switch errorCode {
	case 400:
		var httpErr gophercloud.ErrDefault400
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 401:
		var httpErr gophercloud.ErrDefault401
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 403:
		var httpErr gophercloud.ErrDefault403
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 404:
		var httpErr gophercloud.ErrDefault404
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 405:
		var httpErr gophercloud.ErrDefault405
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 408:
		var httpErr gophercloud.ErrDefault408
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 409:
		var httpErr gophercloud.ErrDefault409
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 429:
		var httpErr gophercloud.ErrDefault429
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 500:
		var httpErr gophercloud.ErrDefault500
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	case 503:
		var httpErr gophercloud.ErrDefault503
		if errors.As(err, &httpErr) {
			return &httpErr.ErrUnexpectedResponseCode, true
		}
		return nil, false
	}

	var unknownErr gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &unknownErr) {
		if unknownErr.Actual == errorCode {
			return &unknownErr, true
		}
	}

	return nil, false
}

func IsNotFound(err error) bool {
	return Is(err, 404)
}

func Any(err error, errorCodes []int) bool {
	for _, c := range errorCodes {
		if Is(err, c) {
			return true
		}
	}
	return false
}

func ErrorWithRequestID(err error, requestID string) error {
	if err == nil {
		return nil
	}
	if requestID == "" {
		return err
	}

	return fmt.Errorf("%w\nRequest ID: %s", err, requestID)
}

func Retry(retryFunc func() error, errorCodes []int, retryCount int, retryDelay time.Duration) error {
	var err error
	for i := 0; i < retryCount; i++ {
		err = retryFunc()
		if err == nil {
			return nil
		}

		needsRetry := false
		for _, errorCode := range errorCodes {
			if Is(err, errorCode) {
				needsRetry = true
				break
			}
		}
		if !needsRetry {
			return err
		}
		time.Sleep(retryDelay)
	}
	return err
}
