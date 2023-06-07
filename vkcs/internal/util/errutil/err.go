package errutil

import (
	"errors"

	"github.com/gophercloud/gophercloud"
)

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
	return false
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
