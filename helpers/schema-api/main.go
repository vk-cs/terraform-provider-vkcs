package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vk-cs/terraform-provider-vkcs/helpers/changelog"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
)

func main() {
	f := flag.NewFlagSet("provider-json", flag.ExitOnError)

	exportSchema := f.String("export", "", "export the schema to the given path/filename. Intended for use in the release process")
	providerName := f.String("provider-name", "vkcs", "the name of the provider")
	changelogPath := f.String("changelog", "CHANGELOG.md", "path to changelog file")

	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error parsing args: %+v", err)
		os.Exit(1)
	}

	data := loadData()

	if exportSchema != nil {
		log.Printf("dumping schema for '%s'", *providerName)
		wrappedProvider := &providerjson.ProviderWrapper{
			ProviderName:  *providerName,
			SchemaVersion: "1",
		}
		cl, err := changelog.NewChangelogFromFile(*changelogPath)
		if err != nil {
			panic(err)
		}
		curVersion := cl.Versions[0].Version
		baseProvider, _ := providerjson.ReadWithWrapper(*exportSchema)
		if err := providerjson.WriteWithWrapper(baseProvider, wrappedProvider, data, *exportSchema, curVersion); err != nil {
			log.Fatalf("error writing provider schema for %q to %q: %+v", *providerName, *exportSchema, err)
		}
	}
}

func loadData() *providerjson.ProviderJSON {
	p := provider.ProviderBase()
	return (*providerjson.ProviderJSON)(p)
}
