package model

import (
	"encoding/json"

	"github.com/rackspace/gophercloud"
)

type Clouds struct {
	Providers []*Provider `hcl:"provider"`
}

type WalkProvidersFn func(provider *Provider, client *gophercloud.ProviderClient) []error

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
		return provider.WalkRegions(client, func(region *Region, client *gophercloud.ProviderClient) (errs []error) {
			return region.Update(client)
		})
	})
}

func (c *Clouds) walk(fn WalkProvidersFn) (errs []error) {
	for _, provider := range c.Providers {
		client, err := provider.Authenticate()
		if err != nil {
			return append(errs, err)
		}
		if werrs := fn(provider, client); werrs != nil {
			errs = append(errs, werrs...)
		}
	}
	return
}
