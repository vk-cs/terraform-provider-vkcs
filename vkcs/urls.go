package vkcs

import "strings"

func baseURL(c ContainerClient, api string) string {
	return c.ServiceURL(api)
}

func getURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id)
}

func instanceActionURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "action")
}

func rootUserURL(c databaseClient, api string, id string) string {
	return c.ServiceURL(api, id, "root")
}

func userDatabasesURL(c databaseClient, api string, id string, userName string) string {
	return c.ServiceURL(api, id, "users", userName, "databases")
}

func userDatabaseURL(c databaseClient, api string, id string, userName string, databaseName string) string {
	return c.ServiceURL(api, id, "users", userName, "databases", databaseName)
}

func userURL(c databaseClient, api string, id string, userName string) string {
	return c.ServiceURL(api, id, "users", userName)
}

func instanceDatabasesURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "databases")
}

func instanceUsersURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "users")
}

func instanceUserURL(c ContainerClient, api string, id string, userName string) string {
	return c.ServiceURL(api, id, "users", userName)
}

func instanceDatabaseURL(c ContainerClient, api string, id string, databaseName string) string {
	return c.ServiceURL(api, id, "databases", databaseName)
}

func instanceCapabilitiesURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "capabilities")
}

func backupScheduleURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "backup_schedule")
}

func deleteURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id)
}

func kubeConfigURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "kube_config")
}

func actionsURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "actions")
}

func upgradeURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "actions", "upgrade")
}

func scaleURL(c ContainerClient, api string, id string) string {
	return c.ServiceURL(api, id, "actions", "scale")
}

func datastoresURL(c ContainerClient, api string) string {
	return c.ServiceURL(api)
}

func datastoreURL(c ContainerClient, api string, dsID string) string {
	return c.ServiceURL(api, dsID)
}

func datastoreParametersURL(c ContainerClient, api string, dsType string, dsVersion string) string {
	return c.ServiceURL(api, dsType, "versions", dsVersion, "parameters")
}

func datastoreCapabilitiesURL(c ContainerClient, api string, dsType string, dsVersion string) string {
	return c.ServiceURL(api, dsType, "versions", dsVersion, "capabilities")
}

func zonesURL(c ContainerClient, api string) string {
	return c.ServiceURL(api, "")
}

func recordsURL(c ContainerClient, api string, zoneID string, recordType string) string {
	return c.ServiceURL(api, zoneID, strings.ToLower(recordType), "")
}

func recordURL(c ContainerClient, api string, zoneID string, recordType string, id string) string {
	return c.ServiceURL(api, zoneID, strings.ToLower(recordType), id)
}
