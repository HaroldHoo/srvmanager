package main

import (
	"github.com/HaroldHoo/srvmanager"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	m := srvmanager.New()

	middleware.DefaultLoggerConfig.Output = m.GetAccessLogWriter()
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("%s\n", "log test ----")
		return c.String(http.StatusOK, "Welcome !")
	})

	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}
