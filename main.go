package main

import (
	"github.com/xlucas/ansible-openstack-inventory/pkg/cmd"
	"github.com/xlucas/ansible-openstack-inventory/pkg/util"
)

func main() {
	if err := cmd.Start(); err != nil {
		util.Die(err.Error(), nil)
	}
}
