package parser

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

var SpreedsheetData GoogleSpreadsheet = GoogleSpreadsheet{
	SpreadsheetId: "1iFFqsu_xruvVKzXAadAAlDBpIuU51v-pfIEU5HeGa8w",
}

func ParseBankData(spreadsheet GoogleSpreadsheet) ([]BankData, error) {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv", SpreedsheetData.SpreadsheetId)

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get data from google sheet: %v", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	rows, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("failed to parse csv data: %v", err)
	}

	var bankData []BankData

	for i, row := range rows {
		if i == 0 {
			continue
		}
		bank := BankData{
			CountryISO2Code: row[0],
			SwiftCode:       row[1],
			CodeType:        row[2],
			Name:            row[3],
			Address:         row[4],
			TownName:        row[5],
			CountryName:     row[6],
			Timezone:        row[7],
		}
		bankData = append(bankData, bank)
	}

	return bankData, nil
}
