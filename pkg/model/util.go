package model

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func GetAnsibleHost(srv servers.Server) string {
	// Return access IP if set
	if srv.AccessIPv4 != "" {
		return srv.AccessIPv4
	}
	if srv.AccessIPv6 != "" {
		return srv.AccessIPv6
	}

	// Find IP from available networks
	var (
		ipv4 string
		ipv6 string
	)
	for _, v := range srv.Addresses {
		for _, v := range v.([]interface{}) {
			v := v.(map[string]interface{})
			if v["OS-EXT-IPS:type"] == "fixed" {
				switch v["version"].(float64) {
				case 4:
					ipv4 = v["addr"].(string)
				case 6:
					ipv6 = v["addr"].(string)
				}
			}
		}
	}

	// Return IPv4 prior to IPv6
	if ipv4 != "" {
		return ipv4
	}
	if ipv6 != "" {
		return ipv6
	}

	// Fallback to DNS
	return srv.Name
}

func GetAnsibleUser(opts *Options, images map[string]images.Image, srv servers.Server) string {
	imageID := srv.Image["id"].(string)
	image := images[imageID]
	if v, ok := image.Metadata[opts.Meta.User]; ok {
		return v
	}
	return opts.FallBackUser
}
