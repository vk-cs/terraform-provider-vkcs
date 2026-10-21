package kms

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/require"
)

func TestWaitForKMSOperation(t *testing.T) {
	t.Run("unfinished operation", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			calls := 0
			var result string
			err := waitForKMSOperation(context.Background(), func() (bool, error) {
				calls++
				if calls == 1 {
					return false, gophercloud.ErrDefault400{ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{Actual: 400, Body: []byte(`{"errors":["unfinished operation"]}`)}}
				}
				result = "ready"
				return true, nil
			})
			require.NoError(t, err)
			require.Equal(t, "ready", result)
			require.Equal(t, 2, calls)
		})
	})
	t.Run("permanent error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			expected := errors.New("permission denied")
			calls := 0
			err := waitForKMSOperation(context.Background(), func() (bool, error) {
				calls++
				return false, expected
			})
			require.ErrorIs(t, err, expected)
			require.Equal(t, 1, calls)
		})
	})
	t.Run("canceled", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			err := waitForKMSOperation(ctx, func() (bool, error) {
				cancel()
				return false, gophercloud.ErrDefault400{ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{Actual: 400, Body: []byte("unfinished operation")}}
			})
			require.ErrorIs(t, err, context.Canceled)
		})
	})
}

func TestKMSStateChangeConfTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		conf := kmsStateChangeConf(func() (interface{}, string, error) {
			return false, kmsStatePending, nil
		}, kmsStateReady)
		conf.Timeout = 10 * time.Second
		_, err := conf.WaitForStateContext(context.Background())
		var timeoutErr *retry.TimeoutError
		require.ErrorAs(t, err, &timeoutErr)
	})
}
