package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"
)

var client *http.Client

func main() {

	client = &http.Client{
		Timeout: time.Second * 5,
	}

	go server()
	go checkDaemon()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
