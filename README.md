# gsloc

`gsloc` is a tool to generate localization files from Google Spreadsheets.

## Usage

```sh
gsloc [command]
```

## Available Commands

### `completion`

Generate the autocompletion script for the specified shell.

### `gen-config-example`

Generate an example config file.

### `gen-loc`

Generate localization files from Google Spreadsheets.

### `help`

Help about any command.

## Flags

```sh
--config string   Config file (default near binary .gsloc.yaml)
-h, --help        Help for gsloc
```

## More Information

Use `gsloc [command] --help` for more information about a command.

---

## Example:

To generate localization files from a Google Spreadsheet, you would use the following command:

```sh
gsloc gen-loc --config path/to/config.yaml
```

To generate an example config file:

```sh
gsloc gen-config-example
```

**Google spreadsheet Example:**

|   | A     | B     | C       |
|---|-------|-------|---------|
| 1 | key   | en    | fr      |
| 2 | hello | Hello | Bonjour |
| 3 | world | World | Monde   |

**Config file for spreadsheet Example:**

```yaml
spreadsheet_id: "***"
sheet_name: "Sheet1"
output_dir: "path/to/output"
rows_to_skip: 1
keys_column: A
locales:
  - column: B
    locale: en
  - column: C
    locale: fr
service_account_json: service_account.json
is_flat: true
```

***How to get the spreadsheet id:***

The spreadsheet id is the long string of characters in the URL of the spreadsheet. For example, in the URL

`https://docs.google.com/spreadsheets/d/1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7/edit#gid=0`

the spreadsheet id is `1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7`.

***How to get the service account json:***

1. Go to the Google Cloud Console. 
2. Select the project that you want to use.
3. In the sidebar on the left, expand APIs & Services and click on Credentials. 
   - Now you have to be here: https://console.cloud.google.com/apis/credentials
4. Click on Create credentials and select Service account.
5. Fill in the required fields and click on Create.
6. Click on Done.

After creating the service account, click on the service account email and then click on Add key. Select JSON and click on Create. The JSON file will be downloaded to your computer.

To get access to the private Google Spreadsheet, you have to share it with the email of the service account.

---