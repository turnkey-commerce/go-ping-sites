package viewmodels

import (
	"strconv"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// ReportViewModel holds the view information for the report.gohtml template
type ReportViewModel struct {
	Title          string
	MonthlyData    map[string]MonthlytItems
	YtdAvgResponse map[string]string
	YtdAvgUptime   map[string]string
	Nav            NavViewModel
	Messages       []string
}

//MonthlytItems defines the slice of items that represent a month of report data
type MonthlytItems []ReportItemViewModel

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

	var monthlyItems map[string]MonthlytItems
	var ytdAvgUptime map[string]string
	var ytdAvgResponse map[string]string
	monthlyItems = make(map[string]MonthlytItems)
	ytdAvgUptime = make(map[string]string)
	ytdAvgResponse = make(map[string]string)
	// Populate the report items.
	for k, v := range reportData {
		monthlyItems[k] = make([]ReportItemViewModel, 12, 12)
		// Traverse twice to get sums first for weighted averages
		totPingsUp := 0
		totPingsDown := 0
		weightedResponseAvg := 0.0
		for _, report := range v {
			if report.PingsUp+report.PingsDown > 0 {
				totPingsUp += report.PingsUp
				totPingsDown += report.PingsDown
			}
		}
		for i, report := range v {
			if report.PingsUp+report.PingsDown > 0 {
				avgResponseStr := strconv.FormatFloat(report.AvgResponse, 'f', 1, 64)
				weightedResponseAvg += report.AvgResponse * float64(report.PingsUp) / float64(totPingsUp)
				uptimePercent := 100.0 * float64(report.PingsUp) / (float64(report.PingsUp) + float64(report.PingsDown))
				uptimePercentStr := strconv.FormatFloat(uptimePercent, 'f', 3, 64)
				monthlyItems[k][i] = ReportItemViewModel{AvgResponse: avgResponseStr,
					UptimePercent: uptimePercentStr}
			} else {
				monthlyItems[k][i] = ReportItemViewModel{AvgResponse: "-",
					UptimePercent: "-"}
			}
		}
		// Calculate ytd.
		uptimePercent := 100.0 * float64(totPingsUp) / (float64(totPingsUp) + float64(totPingsDown))
		ytdAvgResponse[k] = strconv.FormatFloat(weightedResponseAvg, 'f', 1, 64)
		ytdAvgUptime[k] = strconv.FormatFloat(uptimePercent, 'f', 3, 64)
	}

	result := ReportViewModel{
		Title:          "Go Ping Sites - Reports",
		Nav:            nav,
		MonthlyData:    monthlyItems,
		YtdAvgResponse: ytdAvgResponse,
		YtdAvgUptime:   ytdAvgUptime,
		Messages:       messages,
	}

	return result
}
