package monitoring

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ShellScript      string = "#!/bin/bash"
	PowerShellScript string = "#ps1"
	CloudConfig      string = "#cloud-config"
)

var userDataSupportedFormats = []string{ShellScript, CloudConfig, PowerShellScript}

const (
	scriptPermission        = "0777"
	monitoringScriptPath    = "/run/scripts/cloud-monitoring-script.sh"
	userScriptPath          = "/run/scripts/user-script.sh"
	monitoringScriptWinPath = "C:\\Program Files\\Cloudbase Solutions\\Cloudbase-Init\\LocalScripts\\cloud-monitoring-script.ps1"
	userScriptWinPath       = "C:\\Program Files\\Cloudbase Solutions\\Cloudbase-Init\\LocalScripts\\user-script.ps1"
)

func MergeConfigs(monitoringScript string, userData string) (string, error) {
	userData = preprocessScript(userData)
	monitoringScript = preprocessScript(monitoringScript)
	if userData == "" {
		return monitoringScript, nil
	}
	if monitoringScript == "" {
		return userData, nil
	}

	userDataType, header := getScriptType(userData)
	if len(userDataType) == 0 {
		return "", fmt.Errorf("only %s user_data formats are supported, when cloud monitoring is used, given: %s",
			strings.Join(userDataSupportedFormats, ", "), header)
	}

	monitoringScriptType, header := getScriptType(monitoringScript)
	if len(monitoringScriptType) == 0 {
		return "", fmt.Errorf("unknown monitoring script format: %s. Help: set this attribute with `vkcs_cloud_monitoring.script`", header)
	}
	if monitoringScriptType == CloudConfig {
		return "", fmt.Errorf("%s format for monitoring script is not supported. Help: set this attribute with `vkcs_cloud_monitoring.script`", CloudConfig)
	}

	if monitoringScriptType == PowerShellScript {
		switch userDataType {
		case CloudConfig:
			res, err := mergeWinWithCloudConfig(monitoringScript, userData)
			if err != nil {
				return "", fmt.Errorf("failed to merge user_data in cloud-config format with windows monitoring script: %s", err)
			}

			return res, nil
		case PowerShellScript:
			return mergeWinWithPowerShell(monitoringScript, userData)
		case ShellScript:
			return "", fmt.Errorf("monitoring script has %s format, but user_data has %s format. Windows does not have native support for the %s, "+
				"try to rewrite your script in powershell format", PowerShellScript, ShellScript, ShellScript)
		}
	}

	switch userDataType {
	case ShellScript:
		return mergeWithShell(monitoringScript, userData)
	case CloudConfig:
		res, err := mergeWithCloudConfig(monitoringScript, userData)
		if err != nil {
			return "", fmt.Errorf("failed to merge user_data in cloud-config format with monitoring script: %s", err)
		}

		return res, nil
	case PowerShellScript:
		return "", fmt.Errorf("monitoring script has %s format, but user_data has %s format. Unix does not have native support for the %s, "+
			"try to rewrite your script in bash format", ShellScript, PowerShellScript, PowerShellScript)
	}

	return "", fmt.Errorf("unsupported user_data type: %s", userDataType)
}

func preprocessScript(userData string) string {
	return strings.TrimSpace(userData)
}

func getScriptType(script string) (string, string) {
	header, _, _ := strings.Cut(script, "\n")
	header = strings.TrimSpace(header)
	for _, format := range userDataSupportedFormats {
		if header == format {
			return format, ""
		}
	}

	return "", header
}

func mergeWithShell(monitoringScript string, userScript string) (string, error) {
	resultMap := map[string]any{
		"write_files": []map[string]string{
			{
				"path":        userScriptPath,
				"permissions": scriptPermission,
				"content":     userScript,
			},
			{
				"path":        monitoringScriptPath,
				"permissions": scriptPermission,
				"content":     monitoringScript,
			},
		},
		"runcmd": [][]string{
			{"bash", userScriptPath},
			{"bash", monitoringScriptPath},
		},
	}

	return marshalCloudConfig(resultMap)
}

func mergeWinWithPowerShell(monitoringScript string, userScript string) (string, error) {
	resultMap := map[string]any{
		"write_files": []map[string]string{
			{
				"path":        userScriptWinPath,
				"permissions": scriptPermission,
				"content":     userScript,
			},
			{
				"path":        monitoringScriptWinPath,
				"permissions": scriptPermission,
				"content":     monitoringScript,
			},
		},
	}

	return marshalCloudConfig(resultMap)
}

func mergeWithCloudConfig(monitoringScript string, userCloudConfig string) (string, error) {
	var config map[string]any
	err := yaml.Unmarshal([]byte(userCloudConfig), &config)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal cloud-config: %s", err)
	}

	monitoringFile := map[string]string{
		"path":        monitoringScriptPath,
		"permissions": scriptPermission,
		"content":     monitoringScript,
	}
	if err = addMonitoringScriptToConfig(config, monitoringFile); err != nil {
		return "", err
	}

	monitoringCommand := []string{"bash", monitoringScriptPath}
	if block, ok := config["runcmd"]; ok {
		commands, ok := block.([]any)
		if !ok {
			return "", fmt.Errorf("runcmd in cloud_config must be a list")
		}

		config["runcmd"] = append(commands, monitoringCommand)
	} else {
		config["runcmd"] = []any{monitoringCommand}
	}

	return marshalCloudConfig(config)
}

func mergeWinWithCloudConfig(monitoringScript string, userCloudConfig string) (string, error) {
	var config map[string]any
	err := yaml.Unmarshal([]byte(userCloudConfig), &config)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshall cloud-config error: %v", err)
	}

	monitoringFile := map[string]string{
		"path":        monitoringScriptWinPath,
		"permissions": scriptPermission,
		"content":     monitoringScript,
	}
	if err = addMonitoringScriptToConfig(config, monitoringFile); err != nil {
		return "", err
	}

	return marshalCloudConfig(config)
}

func addMonitoringScriptToConfig(config map[string]any, monitoringFile map[string]string) error {
	if block, ok := config["write_files"]; ok {
		if files, ok := block.([]any); ok {
			config["write_files"] = append(files, monitoringFile)
		} else {
			return fmt.Errorf("write_files in cloud_config must be a list")
		}
	} else {
		config["write_files"] = []any{monitoringFile}
	}

	return nil
}

func marshalCloudConfig(config map[string]any) (string, error) {
	out, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to merge user_data with monitoring script, due to marshalling error: %v", err)
	}

	return strings.Join([]string{CloudConfig, string(out)}, "\n"), nil
}
