package util

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

// AppendServerPage is used to extract server information from a page of the
// server list and append it to the given global list.
func AppendServerPage(page pagination.Page, list *[]servers.Server) error {
	pageList, err := servers.ExtractServers(page)
	if err != nil {
		return err
	}
	*list = append(*list, pageList...)
	return nil
}

// AppendImagePage is used to extract image information from a page of the
// image list and append it to the given global map.
func AppendImagePage(page pagination.Page, list map[string]images.Image) error {
	pageList, err := images.ExtractImages(page)
	if err != nil {
		return err
	}
	for _, img := range pageList {
		list[img.ID] = img
	}
	return nil
}
