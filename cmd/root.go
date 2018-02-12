package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xlucas/ansible-openstack-inventory/conf"
	"github.com/xlucas/ansible-openstack-inventory/pkg/util"
)

var (
	list bool
)

func Start() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "openstack-ansible-inventory",
		Short: "Openstack dynamic inventory for ansible",
		Run:   rootCmdFunc,
	}

	c.PersistentFlags().BoolVar(&list, "list", false, "Produces the inventory")

	return c
}

func rootCmdFunc(cmd *cobra.Command, args []string) {
	if !list {
		util.Die("invalid arguments", nil)
	}
	clouds, err := conf.ReadClouds()
	if err != nil {
		util.Die("failed to read clouds configuration", err)
	}
	errs := clouds.Refresh()
	if len(errs) > 0 {
		util.PrintErrors(errs)
		util.Die("failed to refresh host list", nil)
	}
	inventory, err := clouds.BuildInventory()
	if err != nil {
		util.Die("failed to build inventory", err)
	}
	fmt.Fprintln(os.Stdout, string(inventory))
}
