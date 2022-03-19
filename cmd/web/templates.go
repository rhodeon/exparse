package main

import (
	"bytes"
	"html/template"
	"net/http"
)

type TemplateData struct {
	Expr   string
	Result float64
	Error  string
}

func renderTemplate(w http.ResponseWriter, files []string, td TemplateData) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	defer buf.WriteTo(w)

	err = ts.Execute(w, td)
	if err != nil {
		serverError(w, err)
		return
	}
}
