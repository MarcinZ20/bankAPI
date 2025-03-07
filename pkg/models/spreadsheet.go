package models

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
