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

func (p *Provider) walk(client *gophercloud.ProviderClient, walkFunc func(region *Region, client *gophercloud.ProviderClient) []error) (errs []error) {
	syncGroup := new(sync.WaitGroup)

	for _, rg := range p.RegionGroups {
		for _, r := range rg.Regions {
			syncGroup.Add(1)
			go func(region *Region) {
				defer syncGroup.Done()
				if err := walkFunc(region, client); err != nil {
					errs = append(errs, err...)
				}
			}(r)
		}
	}

	syncGroup.Wait()

	return
}

func (p Provider) authenticate() (*gophercloud.ProviderClient, error) {
	return openstack.AuthenticatedClient(p.Identity.getAuthOpts())
}
