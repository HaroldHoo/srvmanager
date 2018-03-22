package main

import (
	"github.com/HaroldHoo/srvmanager"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	m := srvmanager.New()

	gin.DefaultWriter = m.GetAccessLogWriter()
	gin.DefaultErrorWriter = m.GetErrorLogWriter()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("%s\n", "log test ----")
		c.String(http.StatusOK, "Welcome !")
	})

	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}
