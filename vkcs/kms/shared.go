package kms

import (
	"fmt"
	"time"
)

const (
	retriesCount             = 10
	defaultThreshold         = time.Second
	maxItemsInListForIDBuild = 10
)

const (
	EntityFieldID = "id"
)

func entityName(providerTypeName string, entityName string) string {
	return fmt.Sprintf("%s_%s", providerTypeName, entityName)
}
