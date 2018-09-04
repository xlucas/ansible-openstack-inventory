package model

import (
	"encoding/json"

	"github.com/rackspace/gophercloud"
)

type Clouds struct {
	Providers []*Provider `hcl:"provider"`
}

func (c *Clouds) walk(walkFunc func(provider *Provider, client *gophercloud.ProviderClient) []error) (errs []error) {
	for _, provider := range c.Providers {
		client, err := provider.authenticate()
		if err != nil {
			return append(errs, err)
		}
		if werrs := walkFunc(provider, client); werrs != nil {
			errs = append(errs, werrs...)
		}
	}
	return
}

// BuildInventory constructs the ansible inventory and returns it as raw json
// bytes.
func (c *Clouds) BuildInventory(targetEnv string) ([]byte, error) {
	inventory := make(map[string]interface{})

	for _, p := range c.Providers {
		for _, rg := range p.RegionGroups {
			for _, r := range rg.Regions {
				for _, i := range r.instances {
					env, ok := i.Metadata[p.Options.Meta.Env]
					if !ok && targetEnv != "" ||
						ok && env != targetEnv {
						continue
					}
					addToGroups(p, rg, r, i, inventory)
					addToVars(p, rg, r, i, inventory)
				}
			}
		}
	}

	return json.MarshalIndent(inventory, "", "  ")
}

// Refresh updates the Clouds structure with information retrieved from cloud
// providers.
func (c *Clouds) Refresh() []error {
	return c.walk(func(provider *Provider, client *gophercloud.ProviderClient) []error {
		return provider.walk(client, func(region *Region, client *gophercloud.ProviderClient) (errs []error) {
			return region.update(client)
		})
	})
}
