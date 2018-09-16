package model

import (
	"testing"

	"github.com/rackspace/gophercloud"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthOptionsIdentityV2(t *testing.T) {
	identity := Identity{
		Endpoint: "https://keystone.acme.com",
		Username: "BpcRBmKbj1gY",
		Password: "vxWNRLiagH8aEjxA",
		TenantID: "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd",
		Version:  2,
	}
	expected := gophercloud.AuthOptions{
		IdentityEndpoint: identity.Endpoint + "/v2.0",
		Username:         identity.Username,
		Password:         identity.Password,
		TenantID:         identity.TenantID,
		AllowReauth:      true,
	}
	assert.EqualValues(t, expected, identity.GetAuthOpts())
}

func TestGetAuthOptionsIdentityV3(t *testing.T) {
	identity := Identity{
		Endpoint:   "https://keystone.acme.com",
		Username:   "BpcRBmKbj1gY",
		Password:   "vxWNRLiagH8aEjxA",
		ProjectID:  "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd",
		DomainName: "acme",
		Version:    3,
	}
	expected := gophercloud.AuthOptions{
		IdentityEndpoint: identity.Endpoint + "/v3",
		Username:         identity.Username,
		Password:         identity.Password,
		TenantID:         identity.ProjectID,
		DomainName:       identity.DomainName,
		AllowReauth:      true,
	}
	assert.EqualValues(t, expected, identity.GetAuthOpts())
}

func TestGetAuthOptionsIdentityUnsupported(t *testing.T) {
	identity := Identity{
		Endpoint: "https://keystone.acme.com",
		Username: "BpcRBmKbj1gY",
		Password: "vxWNRLiagH8aEjxA",
		TenantID: "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd",
		Version:  1,
	}
	assert.Panics(t, func() { identity.GetAuthOpts() })
}
