package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"bou.ke/monkey"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func monkeyPatchGopherCloud() {
	monkey.Patch(openstack.AuthenticatedClient, func(opts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
		client := &gophercloud.ProviderClient{
			IdentityEndpoint: "https://keystone.acme.com",
			TokenID:          "N28zrJhNSPSAs6piah5JF3dNXaybANi2",
		}
		return client, nil
	})
	monkey.Patch(openstack.NewComputeV2, func(client *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
		return &gophercloud.ServiceClient{}, nil
	})
}

func monkeyPatchImageFetching() {
	t := reflect.TypeOf(&Region{})
	monkey.PatchInstanceMethod(t, "FetchImages", func(r *Region, compute *gophercloud.ServiceClient) error {
		r.images = map[string]images.Image{
			"7b04eb30-f468-4da4-92a9-25d93a1914c1": {},
			"20c5dc91-5a62-4fc2-a122-eeadaadfdf49": {
				Name: "CoreOS",
				Metadata: map[string]string{
					"image_original_user": "core",
				},
			},
		}
		return nil
	})
}

func monkeyPatchServerFetching() {
	t := reflect.TypeOf(&Region{})
	monkey.PatchInstanceMethod(t, "FetchInstances", func(r *Region, compute *gophercloud.ServiceClient) error {
		r.instances = []servers.Server{
			{
				Name:       "web-1",
				AccessIPv4: "130.155.5.7",
				ID:         "0f47385f-2be6-426c-b45f-5b05db68dd11",
				Image: map[string]interface{}{
					"id": "20c5dc91-5a62-4fc2-a122-eeadaadfdf49",
				},
				Metadata: map[string]interface{}{
					"ansible_environment": "production",
					"ansible_groups":      "hardened,web",
					"ansible_hostvar_ssl": "true",
				},
			},
			{
				Name:       "web-2",
				ID:         "f9c33aae-e54a-4ca7-96f8-167f990fd75e",
				AccessIPv6: "481b:a820:6afip1:fa86:b904:88d9:9a0b:9faf",
				Image: map[string]interface{}{
					"id": "7b04eb30-f468-4da4-92a9-25d93a1914c1",
				},
				Metadata: map[string]interface{}{
					"ansible_environment": "production",
					"ansible_groups":      "hardened,web",
					"ansible_hostvar_ssl": "true",
				},
			},
		}
		return nil
	})

}

func TestBuildInventory(t *testing.T) {
	monkeyPatchGopherCloud()
	monkeyPatchImageFetching()
	monkeyPatchServerFetching()

	clouds := &Clouds{
		Providers: []*Provider{
			{
				Name: "acme",
				Identity: Identity{
					Endpoint: "https://keystone.acme.com",
					Username: "BpcRBmKbj1gY",
					Password: "vxWNRLiagH8aEjxA",
					TenantID: "sfgAc5sN3LZUhm2Uho8Sreo0qbUPq8Cd",
					Version:  2,
				},
				Options: Options{
					Meta: Meta{
						Env:            "ansible_environment",
						Groups:         "ansible_groups",
						HostVarsPrefix: "ansible_hostvar_",
						User:           "image_original_user",
					},
					FallBackUser: "admin",
				},
				RegionGroups: []*RegionGroup{
					{
						Name: "eu-east",
						Regions: []*Region{
							{
								Label: "EasternCity",
								Name:  "east-1",
							},
						},
					},
				},
			},
		},
	}

	err := clouds.Refresh()
	assert.Empty(t, err)

	bytes, berr := clouds.BuildInventory("production")
	assert.NoError(t, berr)

	expected := `
	{
	  "_meta": {
		"hostvars": {
		  "0f47385f-2be6-426c-b45f-5b05db68dd11": {
			"ansible_host": "130.155.5.7",
			"ansible_user": "core",
			"provider": "acme",
			"region_group": "eu-east",
			"region_label": "EasternCity",
			"region_name": "east-1",
			"ssl": "true"
		  },
		  "f9c33aae-e54a-4ca7-96f8-167f990fd75e": {
			"ansible_host": "481b:a820:6afip1:fa86:b904:88d9:9a0b:9faf",
			"ansible_user": "admin",
			"provider": "acme",
			"region_group": "eu-east",
			"region_label": "EasternCity",
			"region_name": "east-1",
			"ssl": "true"
		  }
		}
	  },
	  "acme": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "acme_hardened": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "acme_web": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "east-1": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "east-1_hardened": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "east-1_web": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "eu-east": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "eu-east_hardened": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "eu-east_web": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "hardened": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ],
	  "web": [
		"0f47385f-2be6-426c-b45f-5b05db68dd11",
		"f9c33aae-e54a-4ca7-96f8-167f990fd75e"
	  ]
	}`

	assert.JSONEq(t, expected, string(bytes))
}
