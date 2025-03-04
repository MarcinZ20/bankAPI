package parser

// BankData represents a bank data to be processed by the parser
type BankData struct {
	CountryISO2Code string `bson:"country_iso2_code"`
	SwiftCode       string `bson:"swift_code"`
	CodeType        string `bson:"code_type"`
	Name            string `bson:"name"`
	Address         string `bson:"address"`
	TownName        string `bson:"town_name"`
	CountryName     string `bson:"country_name"`
	Timezone        string `bson:"timezone"`
}

// GoogleSpreadsheet represents a spreadsheet data according to Google Sheets API:
//
//	https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets#SpreadsheetProperties
type GoogleSpreadsheet struct {
	SpreadsheetId       string `json:"spreadsheetId"`
	Properties          string `json:"properties"`
	Sheets              string `json:"sheets"`
	NamedRanges         string `json:"namedRanges"`
	SpreadsheetUrl      string `json:"spreadsheetUrl"`
	DeveloperMetadata   string `json:"developerMetadata"`
	DataSource          string `json:"dataSource"`
	DataSourceSchedules string `json:"dataSourceSchedules"`
}
