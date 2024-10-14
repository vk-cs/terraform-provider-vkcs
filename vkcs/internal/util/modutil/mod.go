package modutil

import (
	"runtime/debug"
	"strings"
)

func GetDependencyModuleVersion(module string) (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}

	for _, mod := range info.Deps {
		if mod.Path != module {
			continue
		}

		return strings.TrimPrefix(mod.Version, "v"), true
	}

	return "", false
}
