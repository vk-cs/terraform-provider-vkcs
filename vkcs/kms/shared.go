package kms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	retriesCount             = 10
	defaultThreshold         = 5 * time.Second
	maxItemsInListForIDBuild = 10
)

func isOperationUnfinished(err error) bool {
	return errutil.Is(err, 400) && strings.Contains(err.Error(), "unfinished operation")
}

func waitForKMSOperation(ctx context.Context, check func() (bool, error)) error {
	threshold := defaultThreshold
	for range retriesCount {
		if err := ctx.Err(); err != nil {
			return err
		}
		done, err := check()
		if err != nil && !isOperationUnfinished(err) {
			return err
		}
		if done && err == nil {
			return nil
		}
		timer := time.NewTimer(threshold)
		threshold += defaultThreshold
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("KMS operation did not complete after %d attempts", retriesCount)
}
