package parser

type Bank struct {
	CountryISO2Code string `json:"countryISO2Code"`
	SwiftCode       string `json:"swiftCode"`
	CodeType        string `json:"codeType"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	TownName        string `json:"townName"`
	CountryName     string `json:"countryName"`
	Timezone        string `json:"timezone"`
}

// GoogleSpreadsheet represents a spreadsheet data structure according to Google Sheets API:
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
