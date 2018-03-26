package main

import (
	"github.com/HaroldHoo/srvmanager"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"time"
)

func main() {
	m := srvmanager.New()

	router := martini.Classic()
	router.Logger(log.New(m.GetAccessLogWriter(), "[HTTP] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile))

	router.Get("/", func() (int, string) {
		time.Sleep(5 * time.Second)
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("%s\n", "log test ----")
		return http.StatusOK, "Welcome !"
	})

	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}
