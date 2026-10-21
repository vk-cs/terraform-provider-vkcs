package kms

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	defaultOperationTimeout  = 5 * time.Minute
	defaultThreshold         = 5 * time.Second
	maxItemsInListForIDBuild = 10
)

func isOperationUnfinished(err error) bool {
	return errutil.Is(err, 400) && strings.Contains(err.Error(), "unfinished operation")
}

const (
	kmsStatePending = "pending"
	kmsStateReady   = "ready"
	kmsStateMissing = "missing"
)

func kmsStateChangeConf(refresh retry.StateRefreshFunc, target string) *retry.StateChangeConf {
	pending := []string{kmsStatePending, kmsStateMissing}
	if target == kmsStateMissing {
		pending = []string{kmsStatePending, kmsStateReady}
	}
	return &retry.StateChangeConf{
		Pending:      pending,
		Target:       []string{target},
		Refresh:      refresh,
		Timeout:      defaultOperationTimeout,
		PollInterval: defaultThreshold,
	}
}

// Use the same waiter for operations that do not expose an entity to poll.
func waitForKMSOperation(ctx context.Context, check func() (bool, error)) error {
	_, err := kmsStateChangeConf(func() (interface{}, string, error) {
		done, err := check()
		if isOperationUnfinished(err) {
			return false, kmsStatePending, nil
		}
		if err != nil {
			return nil, "", err
		}
		if done {
			return true, kmsStateReady, nil
		}
		return false, kmsStatePending, nil
	}, kmsStateReady).WaitForStateContext(ctx)
	return err
}
