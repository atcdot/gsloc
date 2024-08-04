package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Command struct {
	cmd *cobra.Command

	configFilePath string
	config         *Config
}

func NewCommand() *Command {
	c := &Command{}

	c.initRootCmd()
	cobra.OnInitialize(c.parseConfig)
	c.initAttrs()

	c.addCommands()

	return c
}

func (c *Command) Execute() {
	if err := c.cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
