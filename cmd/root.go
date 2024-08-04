package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) initRootCmd() {
	c.cmd = &cobra.Command{
		Use:   "gsloc",
		Short: "gsloc is a tool to generate localization files from google spreadsheets",
		Long: `gsloc is a tool to generate localization files from google spreadsheets.
                Complete documentation is available at https://github.com/atcdot/gsloc`,
	}
}
