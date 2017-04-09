package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/apexskier/httpauth"
	"github.com/gorilla/mux"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type reportsController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *reportsController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	year64, err := strconv.ParseInt(vars["year"], 10, 32)
	year := int(year64)
	if err != nil {
		year = time.Now().Year()
	}
	// Get the YTD reports from the database
	ytdReport, err := database.GetYTDReports(controller.DB, year)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetReportViewModel(ytdReport, isAuthenticated, user, messages)
	controller.template.Execute(rw, vm)
	return http.StatusOK, nil
}
