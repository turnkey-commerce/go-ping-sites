package controllers

import (
	"database/sql"
	"net/http"
	"html/template"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/apexskier/httpauth"
)

type reportsController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *reportsController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	// Get the YTD reports from the database
	ytdReport, err := database.GetYTDReports(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetReportViewModel(ytdReport, isAuthenticated, user, messages)
	controller.template.Execute(rw, vm)
	return http.StatusOK, nil
}
