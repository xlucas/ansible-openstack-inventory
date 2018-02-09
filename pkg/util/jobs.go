package util

import (
	"sync"

	"github.com/rackspace/gophercloud"
)

// RunJobs executes jobs using a ServiceClient concurrently.
func RunJobs(service *gophercloud.ServiceClient, jobs []func(service *gophercloud.ServiceClient) error) (errs []error) {
	syncGroup := new(sync.WaitGroup)

	for _, job := range jobs {
		syncGroup.Add(1)
		go func(s *gophercloud.ServiceClient, job func(s *gophercloud.ServiceClient) error) {
			defer syncGroup.Done()
			if err := job(service); err != nil {
				errs = append(errs, err)
			}
		}(service, job)
	}

	syncGroup.Wait()

	return
}
