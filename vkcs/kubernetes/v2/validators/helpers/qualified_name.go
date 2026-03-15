package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	dns1123LabelFmt           string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	dns1123SubdomainFmt       string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"
	dns1123SubdomainMaxLength int    = 253
	qualifiedNameMaxLength    int    = 63
)

var (
	dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")
	qualifiedNameRegexp    = regexp.MustCompile("^" + qualifiedNameFmt + "$")
)

func IsQualifiedName(value string) error {
	parts := strings.Split(value, "/")
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			return errors.New("prefix part of key must be non-empty, got empty prefix part")
		} else if err := isDNS1123Subdomain(prefix); err != nil {
			return err
		}
	default:
		return errors.New("must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyPage'), key doesn't meet the requirements")

	}

	nameLength := len(name)
	if nameLength == 0 {
		return errors.New("name part of key must be non-empty, got empty name part")
	} else if nameLength > qualifiedNameMaxLength {
		return fmt.Errorf("name part of key must be no more than %d characters, got %d characters", qualifiedNameMaxLength, nameLength)
	}

	if !qualifiedNameRegexp.MatchString(name) {
		return fmt.Errorf("name part of key must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'MyKey', 'my.key', '123-key'), got invalid name part '%s'", name)
	}
	return nil
}

func isDNS1123Subdomain(name string) error {
	subdomainLength := utf8.RuneCountInString(name)
	if subdomainLength > dns1123SubdomainMaxLength {
		return fmt.Errorf("prefix part of key must be no more than %d characters, got %d characters", dns1123SubdomainMaxLength, subdomainLength)
	}
	if !dns1123SubdomainRegexp.MatchString(name) {
		return fmt.Errorf("prefix part of key must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character, got invalid prefix part '%s'", name)
	}
	return nil
}
