package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (c *Command) addGenLocCmd() {
	c.cmd.AddCommand(&cobra.Command{
		Use:   "gen-loc",
		Short: "Generate localization files from google spreadsheets",
		Long:  `Generate localization files from google spreadsheets`,
		Run: func(_ *cobra.Command, _ []string) {
			err := c.genLoc()
			if err != nil {
				fmt.Println("Can't generate localization files:", err)
				os.Exit(1)
			}
		},
	})
}

func validateConfig(conf *Config) error {
	if conf.SpreadsheetID == "" {
		return errors.New("spreadsheet_id is required")
	}

	if conf.KeysColumn == "" {
		return errors.New("keys_column is required")
	}

	if len(conf.Locales) == 0 {
		return errors.New("locales is required")
	}

	if conf.OutputDir == "" {
		return errors.New("output_dir is required")
	}

	return nil
}

func (c *Command) genLoc() error {
	fmt.Println("Generating localization files...")

	if err := validateConfig(c.config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// create output directory
	err := os.MkdirAll(c.config.OutputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create output directory: %w", err)
	}

	opts := make([]option.ClientOption, 0)
	opts = append(opts, option.WithScopes(sheets.SpreadsheetsReadonlyScope))

	if c.config.ServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsFile(c.config.ServiceAccountJSON))
	}

	ctx := context.Background()

	sheetsService, err := sheets.NewService(ctx, opts...)
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}

	resp, err := sheetsService.Spreadsheets.Values.Get(c.config.SpreadsheetID, c.config.SheetName).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return errors.New("no data found")
	}

	keyColumnIndex := GetColumnIndexByName(c.config.KeysColumn)

	// create locale files
	for _, locale := range c.config.Locales {
		fmt.Println("Generating locale file for", locale.Locale)

		localeFile, err := os.Create(fmt.Sprintf("%s/%s.json", c.config.OutputDir, locale.Locale))
		if err != nil {
			return fmt.Errorf("can't create locale file: %w", err)
		}
		defer localeFile.Close()

		localeColumnIndex := GetColumnIndexByName(locale.Column)

		_, err = localeFile.WriteString("{\n")
		if err != nil {
			return fmt.Errorf("can't write to locale file: %w", err)
		}

		for i, row := range resp.Values {
			if i < c.config.RowsToSkip {
				continue
			}

			_, err = localeFile.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", row[keyColumnIndex].(string), row[localeColumnIndex].(string)))
			if err != nil {
				return fmt.Errorf("can't write to locale file: %w", err)
			}

			if i < len(resp.Values)-1 {
				_, err := localeFile.WriteString(",")
				if err != nil {
					return fmt.Errorf("can't write to locale file: %w", err)
				}
			}

			_, err = localeFile.WriteString("\n")
			if err != nil {
				return fmt.Errorf("can't write to locale file: %w", err)
			}
		}

		_, err = localeFile.WriteString("}\n")
		if err != nil {
			return fmt.Errorf("can't write to locale file: %w", err)
		}
	}

	fmt.Println("Localization files generated successfully!")
	return nil
}

func GetColumnIndexByName(columnName string) int {
	columnIndex := 0
	for i, r := range columnName {
		columnIndex += int(r-'A'+1) * pow(26, len(columnName)-i-1)
	}
	return columnIndex - 1
}

func pow(i int, i2 int) int {
	if i2 == 0 {
		return 1
	}
	return i * pow(i, i2-1)
}
