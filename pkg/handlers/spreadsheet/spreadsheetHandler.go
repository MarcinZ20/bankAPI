package spreadsheet

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

func FetchSpreadsheetData(spreadsheet *models.GoogleSpreadsheet) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv", spreadsheet.SpreadsheetId)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while fetching data from google spreadsheet: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}
