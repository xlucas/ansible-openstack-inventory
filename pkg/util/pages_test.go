package util

import (
	"errors"
	"testing"

	"github.com/rackspace/gophercloud/openstack/compute/v2/images"

	"github.com/stretchr/testify/assert"

	"bou.ke/monkey"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

func TestAppendImagePage(t *testing.T) {
	list := make(map[string]images.Image)
	pageImages := []images.Image{
		{
			Name: "image-1",
			ID:   "b9dad038-6c85-4ae3-867e-bebddfe62bf1",
		},
		{
			Name: "image-2",
			ID:   "e22a445c-5715-4905-aa3b-f7ffefde5ee8",
		},
	}

	defer monkey.Patch(images.ExtractImages, func(pagination.Page) ([]images.Image, error) {
		return pageImages, nil
	}).Unpatch()

	err := AppendImagePage(images.ImagePage{}, list)
	assert.NoError(t, err)
	assert.EqualValues(t, pageImages[0], list["b9dad038-6c85-4ae3-867e-bebddfe62bf1"])
	assert.EqualValues(t, pageImages[1], list["e22a445c-5715-4905-aa3b-f7ffefde5ee8"])
}

func TestAppendImagePageError(t *testing.T) {
	defer monkey.Patch(images.ExtractImages, func(pagination.Page) ([]images.Image, error) {
		return nil, errors.New("an error occured")
	}).Unpatch()

	err := AppendImagePage(images.ImagePage{}, nil)
	assert.Error(t, err)
}

func TestAppendServerPage(t *testing.T) {
	var list []servers.Server
	pageServers := []servers.Server{
		{
			Name: "server-1",
		},
		{
			Name: "server-2",
		},
	}

	defer monkey.Patch(servers.ExtractServers, func(pagination.Page) ([]servers.Server, error) {
		return pageServers, nil
	}).Unpatch()

	err := AppendServerPage(servers.ServerPage{}, &list)
	assert.NoError(t, err)
	assert.EqualValues(t, pageServers, list)
}

func TestAppendServerPageError(t *testing.T) {
	defer monkey.Patch(servers.ExtractServers, func(pagination.Page) ([]servers.Server, error) {
		return nil, errors.New("an error occured")
	}).Unpatch()

	err := AppendServerPage(servers.ServerPage{}, nil)
	assert.Error(t, err)
}
