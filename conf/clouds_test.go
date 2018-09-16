package conf

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"bou.ke/monkey"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xlucas/ansible-openstack-inventory/pkg/model"
)

func monkeyPatchReadFile() *monkey.PatchGuard {
	config := `
	provider "acme" {
	  options {
	    meta {
		  environment     = "ansible_environment"
		  groups          = "ansible_groups"
		  hostvars_prefix = "ansible_hostvar_"
		  user            = "image_original_user"
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
	return monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return []byte(config), nil
	})
}

func TestReadClouds(t *testing.T) {
	defer monkeyPatchReadFile().Unpatch()

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

func TestReadCloudsExpandHomeError(t *testing.T) {
	defer monkey.Patch(homedir.Expand, func(path string) (string, error) {
		return "", errors.New("an error occured")
	}).Unpatch()
	_, err := ReadClouds()
	assert.Error(t, err)
}

func TestReadCloudsFileError(t *testing.T) {
	defer monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return nil, os.ErrNotExist
	}).Unpatch()
	_, err := ReadClouds()
	assert.Error(t, err)
}

func TestReadCloudsInvalidHCL(t *testing.T) {
	defer monkey.Patch(hcl.Decode, func(out interface{}, in string) error {
		return errors.New("an error occured")
	}).Unpatch()
	_, err := ReadClouds()
	assert.Error(t, err)
}
