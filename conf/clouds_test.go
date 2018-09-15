package conf

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xlucas/ansible-openstack-inventory/pkg/model"

	"bou.ke/monkey"
)

func monkeyPatchReadFile() {
	config := `
	provider "acme" {
		options {
		  meta {
			environment     = "ansible_environment"
			user            = "image_original_user"
			groups          = "ansible_groups"
			hostvars_prefix = "ansible_hostvar_"
		  }
	  
		  fallback_user = "admin"
		}
	  
		identity {
		  version   = 2
		  endpoint  = "https://keystone.acme.com"
		  username  = "BpcRBmKbj1gY"
		  password  = "vxWNRLiagH8aEjxA"
		  tenant_id = "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd"
		}
	  
		regions "eu-east" {
		  region "EasternCity" {
			name = "east-1"
		  }
		}
	  
		regions "eu-west" {
		  region "WesternCity" {
			name = "west-1"
		  }
		}

	  }	  
	`
	monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return []byte(config), nil
	})

}

func TestReadClouds(t *testing.T) {
	monkeyPatchReadFile()

	actual, err := ReadClouds()
	expected := &model.Clouds{
		Providers: []*model.Provider{
			{
				Name: "acme",
				Identity: model.Identity{
					Endpoint: "https://keystone.acme.com",
					Username: "BpcRBmKbj1gY",
					Password: "vxWNRLiagH8aEjxA",
					TenantID: "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd",
					Version:  2,
				},
				Options: model.Options{
					Meta: model.Meta{
						Env:            "ansible_environment",
						Groups:         "ansible_groups",
						HostVarsPrefix: "ansible_hostvar_",
						User:           "image_original_user",
					},
					FallBackUser: "admin",
				},
				RegionGroups: []*model.RegionGroup{
					{
						Name: "eu-east",
						Regions: []*model.Region{
							{
								Label: "EasternCity",
								Name:  "east-1",
							},
						},
					},
					{
						Name: "eu-west",
						Regions: []*model.Region{
							{
								Label: "WesternCity",
								Name:  "west-1",
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}
