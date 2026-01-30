package valid

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/textutil"
)

var (
	ErrInvalidClusterName      = errors.New("invalid cluster name")
	ErrInvalidAvailabilityZone = errors.New("invalid availability zone")
)

const (
	dnsSubdomainFmt       string = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	nameErrorMsg          string = "name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"
	dnsSubdomainMaxLength int    = 25
	dnsSubdomainMinLength int    = 3
)

var (
	dnsSubdomainRegexp = regexp.MustCompile("^" + dnsSubdomainFmt + "$")
	validCharsRegex    = regexp.MustCompile(`[^-a-z0-9]`)
	consecutiveHyphens = regexp.MustCompile(`--`) // check for two consecutive hyphens
)

// ClusterName validates name of cluster.
// Value should match the pattern ^[a-zA-Z][a-zA-Z0-9_.-]*$
func ClusterName(name string) error {
	if len(name) == 0 {
		return ErrInvalidClusterName
	}

	if !textutil.IsLetter(rune(name[0])) {
		return ErrInvalidClusterName
	}

	for _, r := range name[1:] {
		if !textutil.IsLetterDigitSymbol(r, '_', '.', '-') {
			return ErrInvalidClusterName
		}
	}

	return nil
}

// ClusterNameV2 validates cluster name according to DNS subdomain rules
// - Must be 3-25 characters long (after trimming)
// - Must consist only of lowercase latin letters (a-z), digits (0-9), and hyphens (-)
// - Must start and end with alphanumeric character (letter or digit)
// - Must not contain consecutive hyphens (--)
// - Returns detailed error messages matching the original validation
func ClusterNameV2(name string) error {
	name = strings.TrimSpace(name)

	// length validation
	if len(name) < dnsSubdomainMinLength || len(name) > dnsSubdomainMaxLength {
		return fmt.Errorf("%w: cluster name length must be no less than %d and no more than %d characters, got %d",
			ErrInvalidClusterName, dnsSubdomainMinLength, dnsSubdomainMaxLength, len(name))
	}

	// check for consecutive hyphens
	if consecutiveHyphens.MatchString(name) {
		return fmt.Errorf("%w: %s, consecutive hyphens (--) are not allowed",
			ErrInvalidClusterName, nameErrorMsg)
	}

	// format validation using regex
	if !dnsSubdomainRegexp.MatchString(name) {
		// get invalid characters for more informative error message
		invalidChars := getInvalidCharsV2(name)
		if len(invalidChars) > 0 {
			return fmt.Errorf("%w: %s, got invalid characters: %v",
				ErrInvalidClusterName, nameErrorMsg, invalidChars)
		}

		// if there are no invalid characters but regex fails,
		// the issue is with first/last character or other condition
		return fmt.Errorf("%w: %s", ErrInvalidClusterName, nameErrorMsg)
	}

	return nil
}

// getInvalidCharsV2 returns unique invalid characters in the name
func getInvalidCharsV2(name string) []string {
	return removeDuplicatesV2(validCharsRegex.FindAllString(name, -1))
}

// removeDuplicatesV2 removes duplicates from string slice
func removeDuplicatesV2(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
