package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
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
	// Get the YTD reports from the database
	ytdReport, err := database.GetYTDReports(controller.DB, time.Now().Year())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetReportViewModel(ytdReport, isAuthenticated, user, messages)
	controller.template.Execute(rw, vm)
	return http.StatusOK, nil
}
