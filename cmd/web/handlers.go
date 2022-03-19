package main

import (
	"github.com/rhodeon/expression-parser/pkg/solver"
	"github.com/rhodeon/prettylog"
	"net/http"
)

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var homePage = []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
	}

	renderTemplate(w, homePage, TemplateData{})
}

func calculateResult(w http.ResponseWriter, r *http.Request) {
	var resultPage = []string{
		"./ui/html/result.page.gohtml",
		"./ui/html/base.layout.gohtml",
	}

	var homePage = []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
	}

	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	form := r.PostForm
	expr := form.Get("expr")
	prettylog.InfoF("Expression: %s", expr)

	_, err = solver.Validate(expr)
	if err != nil {
		prettylog.ErrorLn(err)
		renderTemplate(w, homePage, TemplateData{Expr: expr, Error: err.Error()})
		return
	}

	result, err := solver.Solve(expr)
	prettylog.InfoF("Result: %s", result)

	if err != nil {
		prettylog.ErrorLn(err)
		renderTemplate(w, homePage, TemplateData{Expr: expr, Error: err.Error()})
		return
	}

	renderTemplate(w, resultPage, TemplateData{Expr: expr, Result: result})
}
