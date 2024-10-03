package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var nginxTimeRegex = regexp.MustCompile(`^\d+(ms|s|m|h|d|w|M|y)?$`)

func NginxTime() validator.String {
	return stringvalidator.RegexMatches(nginxTimeRegex, "")
}
