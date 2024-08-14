package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "./conf.yaml"

type Config struct {
	SpreadsheetID      string         `yaml:"spreadsheet_id" mapstructure:"spreadsheet_id"`
	KeysColumn         string         `yaml:"keys_column" mapstructure:"keys_column"`
	Locales            []LocaleColumn `yaml:"locales" mapstructure:"locales"`
	RowsToSkip         int            `yaml:"rows_to_skip" mapstructure:"rows_to_skip"`
	ServiceAccountJSON string         `yaml:"service_account_json" mapstructure:"service_account_json"`
	SheetName          string         `yaml:"sheet_name" mapstructure:"sheet_name"`
	OutputDir          string         `yaml:"output_dir" mapstructure:"output_dir"`
	IsFlat             bool           `yaml:"is_flat" mapstructure:"is_flat"`
}

type LocaleColumn struct {
	Column string `yaml:"column" mapstructure:"column"`
	Locale string `yaml:"locale" mapstructure:"locale"`
}

func (c *Command) parseConfig() {
	if c.configFilePath == "" {
		c.configFilePath = defaultConfigPath
	}

	viper.SetConfigFile(c.configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// set default values
	c.config = &Config{
		SheetName: "Sheet1",
	}
	if err := viper.Unmarshal(c.config); err != nil {
		fmt.Println("can't unmarshal config:", err)
		os.Exit(1)
	}
}

func (c *Command) addInitConfigCmd() {
	c.cmd.AddCommand(&cobra.Command{
		Use:   "gen-config-example",
		Short: "Generate an example config file",
		Run: func(_ *cobra.Command, _ []string) {
			err := generateConfigExample()
			if err != nil {
				fmt.Println("Can't generate config file:", err)
				os.Exit(1)
			}
		},
	})
}

func generateConfigExample() error {
	defaultConfig := &Config{
		SpreadsheetID: "your-spreadsheet-id",
		KeysColumn:    "A",
		Locales: []LocaleColumn{
			{Column: "B", Locale: "en"},
			{Column: "C", Locale: "de"},
		},
		RowsToSkip:         1,
		ServiceAccountJSON: "service-account.json",
		SheetName:          "Sheet1",
		OutputDir:          "./locales",
		IsFlat:             false,
	}

	f, err := os.Create("conf.yaml")
	if err != nil {
		return fmt.Errorf("can't create config file: %w", err)
	}
	defer f.Close()

	b, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("can't marshal config: %w", err)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("can't write config file: %w", err)
	}

	fmt.Println("Config file generated successfully")
	return nil
}
