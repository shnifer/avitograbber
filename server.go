package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func server() {
	router := fasthttprouter.New()
	router.GET("/", index)
	fasthttp.ListenAndServe(":80", router.Handler)
}

func index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	buf, err := disk.Read("form.html")
	if err != nil {
		return
	}
	ctx.Write(buf)
}
