package model

import (
	"fmt"
	"strings"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func addToGroups(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	defaultGroups := []string{p.Name, rg.Name, r.Name}
	definedGroups := getDefinedGroups(p, s)
	for _, group := range expandGroups(defaultGroups, definedGroups) {
		hostAdd(inventory, group, s.ID)
	}
}

func addToVars(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	hostVarsAdd(inventory, p, rg, r, s)
}

func getDefinedGroups(p *Provider, srv servers.Server) []string {
	if metaValue, ok := srv.Metadata[p.Options.Meta.Groups]; ok {
		return strings.Split(metaValue.(string), ",")
	}
	return nil
}

func expandGroups(defaultGroups, definedGroups []string) (groups []string) {
	groups = append(defaultGroups, definedGroups...)
	for _, definedGroup := range definedGroups {
		for _, defaultGroup := range defaultGroups {
			groups = append(groups, fmt.Sprintf("%s_%s", defaultGroup, definedGroup))
		}
	}
	return
}

func initHostVars(inventory map[string]interface{}) {
	inventory["_meta"] = map[string]interface{}{
		"hostvars": map[string]interface{}{},
	}
}

func hostAdd(inventory map[string]interface{}, group, host string) {
	if _, ok := inventory[group]; !ok {
		inventory[group] = []string{host}
	} else {
		inventory[group] = append(inventory[group].([]string), host)
	}
}

func hostVarsAdd(inventory map[string]interface{}, p *Provider, rg *RegionGroup, r *Region, srv servers.Server) {
	if _, ok := inventory["_meta"]; !ok {
		initHostVars(inventory)
	}
	vars := map[string]interface{}{
		"ansible_host": getAnsibleHost(srv),
		"ansible_user": getAnsibleUser(p, r, srv),
		"provider":     p.Name,
		"region_label": r.Label,
		"region_name":  r.Name,
		"region_group": rg.Name,
	}
	for k, v := range srv.Metadata {
		if strings.HasPrefix(k, p.Options.Meta.HostVarsPrefix) {
			vars[strings.TrimPrefix(k, p.Options.Meta.HostVarsPrefix)] = v
		}
	}
	inventory["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})[srv.ID] = vars
}
