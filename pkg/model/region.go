package model

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/xlucas/ansible-openstack-inventory/pkg/util"
)

type RegionGroup struct {
	Name    string    `hcl:",key"`
	Regions []*Region `hcl:"region"`
}

type Region struct {
	Label     string `hcl:",key"`
	Name      string `hcl:"name"`
	images    map[string]images.Image
	instances []servers.Server
}

func newActiveInstanceFilter() *servers.ListOpts {
	return &servers.ListOpts{
		Status: "ACTIVE",
	}
}

func (r *Region) FetchImages(compute *gophercloud.ServiceClient) error {
	var imgs = make(map[string]images.Image)

	err := images.ListDetail(compute, images.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			if perr := util.AppendImagePage(page, imgs); perr != nil {
				return false, perr
			}
			return true, nil
		},
	)
	if err == nil {
		r.images = imgs
	}

	return err
}

func (r *Region) FetchInstances(compute *gophercloud.ServiceClient) error {
	var instances []servers.Server

	err := servers.List(compute, newActiveInstanceFilter()).EachPage(
		func(page pagination.Page) (bool, error) {
			if perr := util.AppendServerPage(page, &instances); perr != nil {
				return false, perr
			}
			return true, nil
		},
	)
	if err == nil {
		r.instances = instances
	}

	return err
}

func (r *Region) Update(client *gophercloud.ProviderClient) (errs []error) {
	c, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: r.Name,
	})
	if err != nil {
		return append(errs, err)
	}
	jobs := []func(service *gophercloud.ServiceClient) error{
		r.FetchImages,
		r.FetchInstances,
	}
	return util.RunJobs(c, jobs)
}
