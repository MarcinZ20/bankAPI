package spreadsheet

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	tests := []struct {
		name          string
		spreadsheet   *models.GoogleSpreadsheet
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name: "Successful fetch",
			spreadsheet: &models.GoogleSpreadsheet{
				SpreadsheetId: "test123",
			},
			mockResponse:  "CountryISO2,SwiftCode,Type,Name,Address,City,Country\nDE,DEUTDEFFXXX,HQ,Deutsche Bank,Frankfurt,Frankfurt,Germany",
			mockStatus:    http.StatusOK,
			expectedError: false,
		},
		{
			name: "Empty spreadsheet ID",
			spreadsheet: &models.GoogleSpreadsheet{
				SpreadsheetId: "",
			},
			mockResponse:  "",
			mockStatus:    http.StatusNotFound,
			expectedError: true,
		},
		{
			name: "Server error",
			spreadsheet: &models.GoogleSpreadsheet{
				SpreadsheetId: "test123",
			},
			mockResponse:  "Internal Server Error",
			mockStatus:    http.StatusInternalServerError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				assert.Equal(t, "GET", r.Method)

				// Verify that the URL contains spreadsheet ID
				assert.Contains(t, r.URL.String(), tt.spreadsheet.SpreadsheetId)

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()
		})
	}
}
