package valid

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrInvalidNodeGroupName = errors.New("invalid node group name")
)

var (
	ngNameBlackList = []string{"master"}
)

func NodeGroupNameV2(ngName string) error {
	ngName = strings.TrimSpace(ngName)

	if checkNgNameBlackList(ngName) {
		return fmt.Errorf("%w: node group's name can't be a word 'master', got '%s'",
			ErrInvalidNodeGroupName, ngName)
	}

	if len(ngName) < dnsSubdomainMinLength || len(ngName) > dnsSubdomainMaxLength {
		return fmt.Errorf("%w: node group's name length must be no less than %d and no more than %d characters, got %d",
			ErrInvalidNodeGroupName, dnsSubdomainMinLength, dnsSubdomainMaxLength, len(ngName))
	}

	if consecutiveHyphens.MatchString(ngName) {
		return fmt.Errorf("%w: %s, consecutive hyphens (--) are not allowed",
			ErrInvalidNodeGroupName, nameErrorMsg)
	}

	if !dnsSubdomainRegexp.MatchString(ngName) {
		invalidChars := getInvalidCharsV2(ngName)
		if len(invalidChars) > 0 {
			return fmt.Errorf("%w: %s, got invalid characters: %v",
				ErrInvalidNodeGroupName, nameErrorMsg, invalidChars)
		}
		return fmt.Errorf("%w: %s", ErrInvalidNodeGroupName, nameErrorMsg)
	}

	return nil
}

func checkNgNameBlackList(ngName string) bool {
	return slices.Contains(ngNameBlackList, ngName)
}
