package main

import (
	"github.com/HaroldHoo/srvmanager"
	"github.com/kataras/iris"
	"log"
	"net/http"
	"time"
)

func main() {
	m := srvmanager.New()

	router := iris.Default()
	router.Logger().SetOutput(m.GetAccessLogWriter())

	router.Get("/", func(c iris.Context) {
		time.Sleep(5 * time.Second)
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("%s\n", "log test ----")
		c.Writef("Welcome !")
	})

	router.Build()
	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}
