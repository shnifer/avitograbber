package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var formTemplate *template.Template

func initTemplate() {
	buf, err := disk.Read("form.html")
	if err != nil {
		return
	}

	fMap := template.FuncMap{
		"parts":   func() map[string][]string { return siteParts.Parts },
		"names":   func(site string) map[string]string { return siteParts.Names[site] },
		"getName": func(site, part string) string { return siteParts.Names[site][part] },
	}

	formTemplate, err = template.New("form").Funcs(fMap).Parse(string(buf))
	if err != nil {
		log.Println("template parse error: ", err)
		return
	}
}

func server() {
	router := fasthttprouter.New()
	router.GET("/", indexHandler)
	router.GET("/delete", deleteHandler)
	router.GET("/add", addHandler)
	fasthttp.ListenAndServe(":80", router.Handler)
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("recovered ", err)
		}
	}()
	ctx.SetContentType("text/html")
	err := formTemplate.Execute(ctx, askList)
	if err != nil {
		log.Println("execute error")
	}
}

func deleteHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("recovered ", err)
		}
	}()
	defer ctx.Redirect("/", http.StatusTemporaryRedirect)

	param := string(ctx.FormValue("del"))
	n, err := strconv.Atoi(param)
	if err != nil {
		log.Println("strange non-int param ", param)
		return
	}
	DeleteAsList(n)
}

func addHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("recovered ", err)
		}
	}()
	defer ctx.Redirect("/", http.StatusTemporaryRedirect)

	site := string(ctx.FormValue("add"))
	part := string(ctx.FormValue("part"))
	search := string(ctx.FormValue("search"))
	minprice, err := strconv.Atoi(string(ctx.FormValue("minprice")))
	if err != nil {
		minprice = 0
	}
	maxprice, err := strconv.Atoi(string(ctx.FormValue("maxprice")))
	if err != nil {
		maxprice = 0
	}
	physOnly := string(ctx.FormValue("physonly")) != ""
	ask, err := NewAsk(site, part, search, minprice, maxprice, physOnly)
	if err != nil {
		log.Println("NewAsk: ", err)
		return
	}

	AppendAskList(ask)
	doCheck()
}
