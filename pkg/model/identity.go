package model

import (
	"fmt"

	"github.com/rackspace/gophercloud"
)

type Identity struct {
	DomainName string `hcl:"domain"`
	Endpoint   string `hcl:"endpoint"`
	Username   string `hcl:"username"`
	Password   string `hcl:"password"`
	ProjectID  string `hcl:"project_id"`
	TenantID   string `hcl:"tenant_id"`
	Version    int    `hcl:"version"`
}

func (i Identity) GetAuthOpts() gophercloud.AuthOptions {
	if i.Version == 2 {
		return gophercloud.AuthOptions{
			IdentityEndpoint: i.Endpoint + "/v2.0",
			Username:         i.Username,
			Password:         i.Password,
			TenantID:         i.TenantID,
			AllowReauth:      true,
		}
	}
	if i.Version == 3 {
		return gophercloud.AuthOptions{
			IdentityEndpoint: i.Endpoint + "/v3",
			Username:         i.Username,
			Password:         i.Password,
			TenantID:         i.ProjectID,
			DomainName:       i.DomainName,
			AllowReauth:      true,
		}
	}

	panic(fmt.Sprintf("unhandled identity version '%d'", i.Version))
}
