package kms

import "time"

const (
	retriesCount             = 10
	defaultThreshold         = time.Second
	maxItemsInListForIDBuild = 10
)

const (
	EntityFieldID = "id"
)
