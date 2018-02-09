package model

import (
	"strings"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func addToGroups(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	defaultGroups := []string{p.Name, rg.Name, r.Name}
	definedGroups := getDefinedGroups(p, s)
	for _, group := range append(defaultGroups, definedGroups...) {
		hostAdd(inventory, group, s.Name)
	}
}

func addToVars(p *Provider, rg *RegionGroup, r *Region, s servers.Server, inventory map[string]interface{}) {
	hostVarsAdd(inventory, p, r, s)
}

func getDefinedGroups(p *Provider, srv servers.Server) []string {
	if metaValue, ok := srv.Metadata[p.Options.Meta.Groups]; ok {
		return strings.Split(metaValue.(string), ",")
	}
	return nil
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

func hostVarsAdd(inventory map[string]interface{}, p *Provider, r *Region, srv servers.Server) {
	if _, ok := inventory["_meta"]; !ok {
		initHostVars(inventory)
	}
	hostvars := inventory["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})
	hostvars[srv.Name] = map[string]interface{}{
		"ansible_host": getAnsibleHost(srv),
		"ansible_user": getAnsibleUser(p, r, srv),
	}
}
