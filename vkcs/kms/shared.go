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
	defaultThreshold         = time.Second
	maxItemsInListForIDBuild = 10
)

const (
	EntityFieldID = "id"
)

func typeName(providerTypeName string, entityName string) string {
	return fmt.Sprintf("%s_%s", providerTypeName, entityName)
}

func waitForKMSOperation(ctx context.Context, check func() (bool, error)) error {
	for attempt := 0; attempt < retriesCount; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		done, err := check()
		if err != nil && !(errutil.Is(err, 400) && strings.Contains(err.Error(), "unfinished operation")) {
			return err
		}
		if done && err == nil {
			return nil
		}
		timer := time.NewTimer(defaultThreshold)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("KMS operation did not complete after %d attempts", retriesCount)
}
