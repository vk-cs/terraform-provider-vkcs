package kms

import (
	"time"
)

const (
	retriesCount             = 10
	defaultThreshold         = 1 * time.Second
	defaultTimeout           = 1 * time.Minute
	maxItemsInListForIDBuild = 10
)

const (
	diagRetrieveErrorTemplate = "Error retrieving %s from resource: field not found"
	diagSetErrorTemplate      = "Error setting %s: %s"
)
