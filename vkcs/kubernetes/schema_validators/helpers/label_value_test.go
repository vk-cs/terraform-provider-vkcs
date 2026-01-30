package helpers

import (
	"strings"
	"testing"
)

func TestIsValidTaintValue(t *testing.T) {
	successCases := []string{
		"simple",
		"now-with-dashes",
		"1-starts-with-num",
		"end-with-num-1",
		"1234",                  // only num
		strings.Repeat("a", 63), // to the limit
		"",                      // empty value
	}
	for i := range successCases {
		if err := IsValidLabelValue(successCases[i]); err != nil {
			t.Errorf("case %s expected success: %v", successCases[i], err)
		}
	}

	errorCases := []string{
		"nospecialchars%^=@",
		"Tama-nui-te-rā.is.Māori.sun",
		"\\backslashes\\are\\bad",
		"-starts-with-dash",
		"ends-with-dash-",
		".starts.with.dot",
		"ends.with.dot.",
		strings.Repeat("a", 64), // over the limit
	}
	for i := range errorCases {
		if err := IsValidLabelValue(errorCases[i]); err == nil {
			t.Errorf("case[%d] expected failure", i)
		}
	}
}
