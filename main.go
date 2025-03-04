package main

import (
	"fmt"

	"github.com/MarcinZ20/bankAPI/internal/parser"
)

func main() {
	bankData, err := parser.ParseBankData(parser.SpreedsheetData)
	if err != nil {
		fmt.Printf("failed to parse bank data: %v\n", err)
		return
	}

	for _, bank := range bankData {
		fmt.Printf("CountryISO2Code: %s\n", bank.CountryISO2Code)
		fmt.Printf("SWIFTCode: %s\n", bank.SwiftCode)
		fmt.Printf("CodeType: %s\n", bank.CodeType)
		fmt.Printf("Name: %s\n", bank.Name)
		fmt.Printf("Address: %s\n", bank.Address)
		fmt.Printf("TownName: %s\n", bank.TownName)
		fmt.Printf("CountryName: %s\n", bank.CountryName)
		fmt.Printf("Timezone: %s\n", bank.Timezone)
		fmt.Println()
	}
}
