package viewmodels

import (
	"strconv"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/httpauth"
)

// ReportViewModel holds the view information for the report.gohtml template
type ReportViewModel struct {
	Title      string
	ReportData map[string]ReportItems
	Nav        NavViewModel
	Messages   []string
}

//ReportItems defines the slice of items that represent a month of report data
type ReportItems []ReportItemViewModel

//ReportItemViewModel defines the individual items that will be reported.
type ReportItemViewModel struct {
	AvgResponse   string
	UptimePercent string
}

// GetReportViewModel populates the items required by the report.gohtml view
func GetReportViewModel(reportData map[string]database.Reports, isAuthenticated bool, user httpauth.UserData, messages []string) ReportViewModel {
	nav := NavViewModel{
		Active:          "reports",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	var reportItems map[string]ReportItems
	reportItems = make(map[string]ReportItems)
	// Populate the report items.
	for k, v := range reportData {
		reportItems[k] = make([]ReportItemViewModel, 12, 12)
		for i, report := range v {
			if report.PingsUp+report.PingsDown > 0 {
				avgResponseStr := strconv.FormatFloat(report.AvgResponse, 'f', 1, 64)
				uptimePercent := 100.0 * float64(report.PingsUp) / (float64(report.PingsUp) + float64(report.PingsDown))
				uptimePercentStr := strconv.FormatFloat(uptimePercent, 'f', 3, 64)
				reportItems[k][i] = ReportItemViewModel{AvgResponse: avgResponseStr,
					UptimePercent: uptimePercentStr}
			} else {
				reportItems[k][i] = ReportItemViewModel{AvgResponse: "-",
					UptimePercent: "-"}
			}
		}
	}

	result := ReportViewModel{
		Title:      "Go Ping Sites - Reports",
		Nav:        nav,
		ReportData: reportItems,
		Messages:   messages,
	}

	return result
}
