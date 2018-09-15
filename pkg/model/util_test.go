package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func TestGetAnsibleHostFromAccessIPv4(t *testing.T) {
	host := "130.155.5.7"
	server := servers.Server{
		AccessIPv4: host,
	}
	assert.Equal(t, host, GetAnsibleHost(server))
}

func TestGetAnsibleHostFromAccessIPv6(t *testing.T) {
	host := "481b:a820:6af1:fa86:b904:88d9:9a0b:9faf"
	server := servers.Server{
		AccessIPv6: host,
	}
	assert.Equal(t, host, GetAnsibleHost(server))
}

func TestGetAnsibleHostFromIPv4Network(t *testing.T) {
	host := "130.155.5.7"
	server := servers.Server{
		Addresses: map[string]interface{}{
			"Ext-Net": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS:type": "fixed",
					"version":         float64(4),
					"addr":            host,
				},
			},
		},
	}
	assert.Equal(t, host, GetAnsibleHost(server))
}

func TestGetAnsibleHostFromIPv6Network(t *testing.T) {
	host := "481b:a820:6afip1:fa86:b904:88d9:9a0b:9faf"
	server := servers.Server{
		Addresses: map[string]interface{}{
			"Ext-Net": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS:type": "fixed",
					"version":         float64(6),
					"addr":            host,
				},
			},
		},
	}
	assert.Equal(t, host, GetAnsibleHost(server))
}

func TestGetAnsibleHostFromName(t *testing.T) {
	host := "server.acme.com"
	server := servers.Server{Name: host}
	assert.Equal(t, host, GetAnsibleHost(server))
}

func TestGetAnsibleUserWithMeta(t *testing.T) {
	opts := &Options{
		Meta: Meta{
			User: "image_original_user",
		},
		FallBackUser: "admin",
	}
	images := map[string]images.Image{
		"20c5dc91-5a62-4fc2-a122-eeadaadfdf49": {
			Metadata: map[string]string{
				"image_original_user": "core",
			},
		},
	}
	server := servers.Server{
		Image: map[string]interface{}{
			"id": "20c5dc91-5a62-4fc2-a122-eeadaadfdf49",
		},
	}
	assert.Equal(t, "core", GetAnsibleUser(opts, images, server))
}

func TestGetAnsibleUserWithoutMeta(t *testing.T) {
	opts := &Options{
		Meta: Meta{
			User: "image_original_user",
		},
		FallBackUser: "admin",
	}
	images := map[string]images.Image{
		"7b04eb30-f468-4da4-92a9-25d93a1914c1": {
			Metadata: make(map[string]string),
		},
	}
	server := servers.Server{
		Image: map[string]interface{}{
			"id": "7b04eb30-f468-4da4-92a9-25d93a1914c1",
		},
	}
	assert.Equal(t, "admin", GetAnsibleUser(opts, images, server))
}
