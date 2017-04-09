package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type reportsController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *reportsController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	vals := req.URL.Query()
	year := time.Now().Year()
	yearParms, ok := vals["year"]
	if ok && len(yearParms) > 0 {
		year64, err := strconv.ParseInt(yearParms[0], 10, 32)
		// Only get the year if there was no error.
		if err == nil {
			year = int(year64)
		}
	}

	// Get the Reporting years from the database
	years, err := database.GetReportYears(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Get the YTD reports from the database
	ytdReport, err := database.GetYTDReports(controller.DB, year)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetReportViewModel(year, years, ytdReport, isAuthenticated, user, messages)
	controller.template.Execute(rw, vm)
	return http.StatusOK, nil
}
