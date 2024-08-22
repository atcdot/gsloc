package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

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
			c.parseConfig()
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

	if conf.ServiceAccountJSON == "" {
		return errors.New("service_account_json is required (See https://developers.google.com/workspace/guides/create-credentials?#service-account)")
	}

	return nil
}

type translation struct {
	key   string
	value string
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
	opts = append(opts, option.WithCredentialsFile(c.config.ServiceAccountJSON))

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

		localeColumnIndex := GetColumnIndexByName(locale.Column)

		rows := make([]translation, 0, len(resp.Values))
		for i, row := range resp.Values {
			if i < c.config.RowsToSkip {
				continue
			}

			// skip empty rows and rows without locale
			if len(row) <= keyColumnIndex || len(row) <= localeColumnIndex {
				continue
			}

			valueStr, ok := row[localeColumnIndex].(string)
			if !ok || valueStr == "" {
				continue
			}

			rows = append(rows, translation{
				key:   row[keyColumnIndex].(string),
				value: valueStr,
			})
		}

		switch c.config.IsFlat {
		case true:
			if err := writeLocaleFileFlat(c.config.OutputDir, locale, rows); err != nil {
				return fmt.Errorf("can't write locale file: %w", err)
			}
		default:
			if err := writeLocaleFileTree(c.config.OutputDir, locale, rows); err != nil {
				return fmt.Errorf("can't write locale file: %w", err)
			}
		}
	}

	fmt.Println("Localization files generated successfully!")
	return nil
}

func writeLocaleFileTree(outputDir string, locale LocaleColumn, rows []translation) error {
	localeFile, err := os.Create(fmt.Sprintf("%s/%s.json", outputDir, locale.Locale))
	if err != nil {
		return fmt.Errorf("can't create locale file: %w", err)
	}
	defer localeFile.Close()

	output := make(map[string]interface{})

	for _, row := range rows {
		set(output, row.key, row.value)
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshal locale file: %w", err)
	}

	_, err = localeFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("can't write to locale file: %w", err)
	}

	return nil
}

func writeLocaleFileFlat(outputDir string, locale LocaleColumn, rows []translation) error {
	localeFile, err := os.Create(fmt.Sprintf("%s/%s.json", outputDir, locale.Locale))
	if err != nil {
		return fmt.Errorf("can't create locale file: %w", err)
	}
	defer localeFile.Close()

	_, err = localeFile.WriteString("{\n")
	if err != nil {
		return fmt.Errorf("can't write to locale file: %w", err)
	}

	i := 0
	for _, tr := range rows {
		_, err = localeFile.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", tr.key, tr.value))
		if err != nil {
			return fmt.Errorf("can't write to locale file: %w", err)
		}

		if i < len(rows)-1 {
			_, err := localeFile.WriteString(",")
			if err != nil {
				return fmt.Errorf("can't write to locale file: %w", err)
			}
		}

		_, err = localeFile.WriteString("\n")
		if err != nil {
			return fmt.Errorf("can't write to locale file: %w", err)
		}

		i++
	}

	_, err = localeFile.WriteString("}\n")
	if err != nil {
		return fmt.Errorf("can't write to locale file: %w", err)
	}
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

// Set function to set nested map values based on dotted keys
func set(m map[string]interface{}, key string, value string) {
	keys := strings.Split(key, ".")
	for i, k := range keys {
		if i == len(keys)-1 {
			m[k] = value
		} else {
			if _, ok := m[k]; !ok {
				m[k] = make(map[string]interface{})
			}
			m = m[k].(map[string]interface{})
		}
	}
}
