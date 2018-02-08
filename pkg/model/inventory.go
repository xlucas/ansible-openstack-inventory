package model

import (
	"strings"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func addToGroups(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	groups := []string{p.Name, rg.Name, r.Name}

	if str, ok := s.Metadata[p.Options.Meta.Groups]; ok {
		userGroups := strings.Split(str.(string), ",")
		groups = append(groups, userGroups...)
	}
	for _, group := range groups {
		if _, ok := inventory[group]; !ok {
			inventory[group] = []string{}
		}
		inventory[group] = append(inventory[group].([]string), s.Name)
	}
}

func addToVars(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	if inventory["_meta"] == nil {
		inventory["_meta"] = map[string]interface{}{
			"hostvars": map[string]interface{}{},
		}
	}
	hostvars := inventory["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})
	hostvars[s.Name] = map[string]interface{}{
		"ansible_host": getIpAddress(s),
		"ansible_user": getSSHUser(s.Image["id"].(string), p, r),
	}
}

func getSSHUser(imageID string, provider *Provider, region *Region) string {
	image := region.images[imageID]
	if v, ok := image.Metadata[provider.Options.Meta.User]; ok {
		return v
	}
	return provider.Options.FallBackUser
}
