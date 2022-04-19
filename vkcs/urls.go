package vkcs

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
