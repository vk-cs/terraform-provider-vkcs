package main

import (
	"flag"
	"log"
	"os"

	"github.com/vk-cs/terraform-provider-vkcs/helpers/changelog"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/json"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	transform "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/transform/provider"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
)

func main() {
	f := flag.NewFlagSet("provider-json", flag.ExitOnError)

	exportSchema := f.String("export", "", "export the schema to the given path/filename. Intended for use in the release process")
	providerName := f.String("provider-name", "vkcs", "the name of the provider")
	changelogPath := f.String("changelog", "CHANGELOG.md", "path to changelog file")

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing args: %+v", err)
	}

	data := loadData()

	if exportSchema != nil {
		log.Printf("dumping schema for '%s'", *providerName)

		cl, err := changelog.NewChangelogFromFile(*changelogPath)
		if err != nil {
			log.Fatalf("error parsing changelog: %+v", err)
		}

		baseProvider, _ := json.ReadWithWrapper(*exportSchema)
		wrappedProvider := &schema.ProviderWrapper{
			ProviderName:    *providerName,
			ProviderVersion: cl.Versions[0].Version,
			SchemaVersion:   "1",
		}

		wrappedProvider, err = transform.WrappedProviderFromRaw(data, baseProvider, wrappedProvider)
		if err != nil {
			log.Fatalf("error transforming provider into json schema: %+v", err)
		}

		if err := json.WriteWithWrapper(wrappedProvider, *exportSchema); err != nil {
			log.Fatalf("error writing provider schema for %q to %q: %+v", *providerName, *exportSchema, err)
		}
	}
}

func loadData() *schema.ProviderJSON {
	return &schema.ProviderJSON{SDKProvider: provider.SDKProviderBase(), Provider: provider.ProviderBase()}
}
