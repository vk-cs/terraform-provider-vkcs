package helpers

import (
	"fmt"
	"regexp"
)

const (
	qnameCharFmt        string = "[A-Za-z0-9]"
	qnameExtCharFmt     string = "[-A-Za-z0-9_.]"
	qualifiedNameFmt    string = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	labelValueFmt       string = "(" + qualifiedNameFmt + ")?"
	labelValueMaxLength int    = 63
)

var labelValueRegexp = regexp.MustCompile("^" + labelValueFmt + "$")

func IsValidLabelValue(label string) error {
	valueLength := len(label)
	if valueLength > labelValueMaxLength {
		return fmt.Errorf("value must be no more than %d characters, got %d characters", labelValueMaxLength, valueLength)
	}
	if !labelValueRegexp.MatchString(label) {
		return fmt.Errorf("value must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue', 'my_value', '12345'), got invalid value '%s'", label)
	}
	return nil
}
