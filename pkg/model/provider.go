package model

import (
	"sync"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type Provider struct {
	Identity     Identity       `hcl:"identity"`
	Options      Options        `hcl:"options"`
	Name         string         `hcl:",key"`
	RegionGroups []*RegionGroup `hcl:"regions"`
}

type WalkRegionsFn func(region *Region, client *gophercloud.ProviderClient) []error

// Authenticate is used to authenticate to the identity endpoint.
func (p Provider) Authenticate() (*gophercloud.ProviderClient, error) {
	return openstack.AuthenticatedClient(p.Identity.GetAuthOpts())
}

// WalkRegions is used to run fn for each region.
func (p *Provider) WalkRegions(client *gophercloud.ProviderClient, fn WalkRegionsFn) (errs []error) {
	syncGroup := new(sync.WaitGroup)

	for _, rg := range p.RegionGroups {
		for _, r := range rg.Regions {
			syncGroup.Add(1)
			go func(region *Region) {
				defer syncGroup.Done()
				if err := fn(region, client); err != nil {
					errs = append(errs, err...)
				}
			}(r)
		}
	}

	syncGroup.Wait()

	return
}
