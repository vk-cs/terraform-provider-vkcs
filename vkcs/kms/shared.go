package kms

import (
	"errors"
	"time"
)

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

var ErrResponseDecodeFail = errors.New("failed to decode response")
