package cmd

func (c *Command) initAttrs() {
	c.cmd.PersistentFlags().StringVar(&c.configFilePath, "config", "", "config file (default near binary .gsloc.yaml)")
}
