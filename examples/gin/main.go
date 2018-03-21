package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/HaroldHoo/srvmanager"
)

func main() {
	m := srvmanager.New()

	gin.DefaultWriter = m.GetAccessLogWriter()
	gin.DefaultErrorWriter = m.GetErrorLogWriter()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}

