package schema_validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*InsecureRegistryValidator)(nil)
)

const (
	minPort = 1
	maxPort = 65535

	registryValidator = "^[a-zA-Z0-9:/.-]+$"
	wildcardPrefix    = "*."
)

var (
	registryRegex = regexp.MustCompile(registryValidator)
)

type InsecureRegistryValidator struct{}

func (v InsecureRegistryValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v InsecureRegistryValidator) MarkdownDescription(ctx context.Context) string {
	return "Insecure registry must be valid URL"
}

//nolint:staticcheck
func (v InsecureRegistryValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	registry := req.ConfigValue.ValueString()
	if err := isValidInsecureRegistry(registry); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid insecure registry URL",
			err.Error(),
		)
	}

}

//nolint:staticcheck
func isValidInsecureRegistry(registry string) error {
	// check the registry for correctness if an wildcard is specified
	if strings.Contains(registry, "*") {
		return isWildcardValid(registry)
	}

	if err := isRegistryValid(registry); err != nil {
		return err
	}

	// split into the port part and the path
	// first, look for the first colon to separate the port
	parts := strings.SplitN(registry, ":", 2)

	// process the main part (before the port)
	mainPart := parts[0]

	// check if there is a path in the main part
	domain, path := splitDomainAndPath(mainPart)

	if err := isDomainValidWithError(domain); err != nil {
		return err
	}

	// if there is a port
	if len(parts) == 2 {
		// process the part with the port and the possible path after it
		portPart := parts[1]

		// separate the port and path (if any)
		portPathParts := strings.SplitN(portPart, "/", 2)
		port := portPathParts[0]

		// check the port
		if err := isPortValidWithError(port); err != nil {
			return err
		}

		// if there is a path after the port
		if len(portPathParts) == 2 {
			portPath := portPathParts[1]

			// if there was a path in the main part, they must match
			if path != "" && path != portPath {
				return fmt.Errorf("Inconsistent path specification: path after port (%s) must match main path (%s)", portPath, path)
			}
			// if there was no path in the main part, keep the path from the port part
			if path == "" {
				path = portPath
			}
		}
	}

	// check the path if it exists
	if path != "" {
		if err := isPathValidWithError(path); err != nil {
			return err
		}
	}

	return nil
}

// splitDomainAndPath splits a string into domain and path
func splitDomainAndPath(s string) (domain, path string) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

//nolint:staticcheck
func isWildcardValid(registry string) error {
	if !strings.HasPrefix(registry, wildcardPrefix) {
		return fmt.Errorf("Wildcard (*) must be at the beginning and followed by a dot, like '*.example.com'")
	}

	if registry == wildcardPrefix {
		return fmt.Errorf("Wildcard prefix '*/' must be followed by a domain name")
	}

	for i := 2; i < len(registry); i++ {
		c := registry[i]
		if c == '/' {
			return fmt.Errorf("Wildcard registry cannot contain paths (found '/')")
		}
		if c == ':' {
			return fmt.Errorf("Wildcard registry cannot contain port specification (found ':')")
		}
		if c == '*' {
			return fmt.Errorf("Wildcard registry cannot contain multiple asterisks (*)")
		}
	}

	return nil
}

//nolint:staticcheck
func isRegistryValid(registry string) error {
	if registry == "" {
		return fmt.Errorf("Registry cannot be empty")
	}

	if !registryRegex.MatchString(registry) {
		return fmt.Errorf("Registry contains invalid characters. Allowed: letters (a-z, A-Z), numbers (0-9), and characters: : / . -")
	}

	length := len(registry)
	for i := range length {
		if registry[i] == '/' {
			if length == i+1 {
				return fmt.Errorf("Registry cannot end with a slash (/)")
			}
			if registry[i+1] == '/' {
				return fmt.Errorf("Registry cannot contain consecutive slashes (//)")
			}
		}
	}

	return nil
}

//nolint:staticcheck
func isPathValidWithError(path string) error {
	// check for empty path segments
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment == "" {
			return fmt.Errorf("Path cannot contain empty segments (consecutive slashes) at position %d", i)
		}
	}
	return nil
}

//nolint:staticcheck
func isDomainValidWithError(domain string) error {
	if domain == "" {
		return fmt.Errorf("Domain/IP part cannot be empty")
	}

	if isIP(domain) {
		return isIPValidWithError(domain)
	}

	hasDot := false
	for i, r := range domain {
		if r == ':' || r == '/' {
			return fmt.Errorf("Domain cannot contain '%c' - port and path should be separated with proper syntax", r)
		}

		if r == '.' {
			hasDot = true
			if i == 0 || i == len(domain)-1 {
				return fmt.Errorf("Domain cannot start or end with a dot (.)")
			}
		}
	}

	if !hasDot {
		return fmt.Errorf("Domain name must contain at least one dot (e.g., 'example.com')")
	}

	return nil
}

func isIP(domain string) bool {
	octets := strings.Split(domain, ".")
	if len(octets) != 4 {
		return false
	}

	for _, octet := range octets {
		for _, digit := range octet {
			if digit < '0' || digit > '9' {
				return false
			}
		}
	}

	return true
}

func isIPValidWithError(domain string) error {
	octets := strings.Split(domain, ".")
	if len(octets) != 4 {
		return fmt.Errorf("IP address must have exactly 4 octets, got %d", len(octets))
	}

	for i, octet := range octets {
		if octet == "" {
			return fmt.Errorf("IP address octet %d is empty", i+1)
		}

		// check for leading zeros
		if len(octet) > 1 && octet[0] == '0' {
			return fmt.Errorf("IP address octet %d cannot have leading zeros", i+1)
		}

		number, err := strconv.Atoi(octet)
		if err != nil {
			return fmt.Errorf("IP address octet %d contains non-numeric characters", i+1)
		}

		if number < 0 || number > 255 {
			return fmt.Errorf("IP address octet %d must be between 0 and 255, got %d", i+1, number)
		}
	}

	return nil
}

//nolint:staticcheck
func isPortValidWithError(portStr string) error {
	if portStr == "" {
		return fmt.Errorf("Port number cannot be empty")
	}

	if portStr[0] == '0' {
		if portStr == "0" {
			return fmt.Errorf("Port 0 is reserved and cannot be used")
		}
		return fmt.Errorf("Port cannot have leading zeros")
	}

	port := 0
	for i, r := range portStr {
		if r < '0' || r > '9' {
			return fmt.Errorf("Port contains non-digit character '%c' at position %d", r, i)
		}

		port = port*10 + int(r-'0')
		if port > maxPort {
			return fmt.Errorf("Port %d exceeds maximum allowed value %d", port, maxPort)
		}
	}

	if port < minPort {
		return fmt.Errorf("Port must be between %d and %d, got %d", minPort, maxPort, port)
	}

	return nil
}
