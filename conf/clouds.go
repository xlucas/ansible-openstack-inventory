package conf

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/xlucas/ansible-openstack-inventory/pkg/model"
)

const (
	// CloudsDefaultPath is the default path for the clouds configuration file.
	CloudsDefaultPath = "~/.clouds.hcl"
)

// ReadClouds parses the content of the file containing user-configured clouds.
func ReadClouds() (*model.Clouds, error) {
	var clouds = new(model.Clouds)

	path, err := homedir.Expand(CloudsDefaultPath)
	if err != nil {
		return nil, err
	}
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = hcl.Decode(&clouds, string(confBytes))
	if err != nil {
		return nil, err
	}

	return clouds, err
}
