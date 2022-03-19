package main

import (
	"github.com/rhodeon/prettylog"
	"net/http"
)

func serverError(w http.ResponseWriter, err error) {
	prettylog.Error(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
