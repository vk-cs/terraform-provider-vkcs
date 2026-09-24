package kms

import "time"

type SecretStatus = string

const (
	SecretStatusActive     SecretStatus = "active"
	SecretStatusInProgress SecretStatus = "in_progress"
)

type KeyStatus = string

const (
	KeyStatusActive     = "active"
	KeyStatusInProgress = "in_progress"
)

const (
	defaultTimeout = 1 * time.Minute
)
