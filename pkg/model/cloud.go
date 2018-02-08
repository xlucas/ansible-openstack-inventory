package model

import (
	"encoding/json"
	"sync"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type Clouds struct {
	Providers []*Provider `hcl:"provider"`
}

func (c *Clouds) walk(walkFunc func(provider *Provider, client *gophercloud.ProviderClient) []error) (errs []error) {
	for _, provider := range c.Providers {
		client, err := provider.authenticate()
		if err != nil {
			errs = append(errs, err)
			return
		}
		if werrs := walkFunc(provider, client); werrs != nil {
			errs = append(errs, werrs...)
		}
	}
	return
}

// BuildInventory constructs the ansible inventory and returns it as raw json
// bytes.
func (c *Clouds) BuildInventory() ([]byte, error) {
	inventory := make(map[string]interface{})

	for _, p := range c.Providers {
		for _, rg := range p.RegionGroups {
			for _, r := range rg.Regions {
				for _, i := range r.instances {
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
			c, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
				Region: region.Name,
			})
			if err != nil {
				errs = append(errs, err)
				return
			}

			jobs := []func(service *gophercloud.ServiceClient) error{
				region.fetchImages,
				region.fetchInstances,
			}

			var group = new(sync.WaitGroup)

			for _, j := range jobs {
				group.Add(1)
				go func(compute *gophercloud.ServiceClient, job func(service *gophercloud.ServiceClient) error) {
					defer group.Done()
					if err := job(compute); err != nil {
						errs = append(errs, err)
					}
				}(c, j)
			}

			group.Wait()
			return
		})
	})
}
