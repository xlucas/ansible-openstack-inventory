package util

import (
	"errors"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/stretchr/testify/assert"
)

func TestRunJobs(t *testing.T) {
	jobs := []func(service *gophercloud.ServiceClient) error{
		func(*gophercloud.ServiceClient) error {
			return nil
		},
	}
	errs := RunJobs(nil, jobs)
	assert.Empty(t, errs)
}

func TestRunJobsError(t *testing.T) {
	jobs := []func(service *gophercloud.ServiceClient) error{
		func(*gophercloud.ServiceClient) error {
			return errors.New("an error occured")
		},
		func(*gophercloud.ServiceClient) error {
			return errors.New("another error occured")
		},
	}
	errs := RunJobs(nil, jobs)
	assert.Len(t, errs, 2)
}
