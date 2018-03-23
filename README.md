# srvmanager
1. graceful shutdown/reload/update reload. （平滑地 关闭/重启/迭代二进制重启）
2. manager process via -s and -pid option. （通过 可选的-pid 和 -s reload 或 kill -HUP pid 像Nginx一样优雅地重启服务）
3. log something easier. （记录日志更便捷）

## Usage
```
./main/bin/server -h

Usage of ./main/bin/server:
  -accesslog string
    	log file (default "/var/log/server_access.log")
  -errorlog string
    	log file (default "/var/log/server_error.log")
  -pid string
    	pid file (default "/var/run/server.pid")
  -d	Start as deamon. (default true)
  -p string
    	Listen port (default "8080")
  -s string
    	(When used with the -pid option, the pid will be from specified pidfile.)
    	Send a signal to the process.  The argument signal can be one of: start stop reload restart,
    	The following table shows the corresponding system signals:
    	stop	SIGTERM
    	reload	SIGHUP
    	restart	SIGHUP
    	 (default "start")
```

---

## Examples
   1. Via https://github.com/gin-gonic/gin : （配合优秀的gin框架）
   > https://github.com/HaroldHoo/srvmanager/tree/master/examples/gin

```
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
		time.Sleep(2 * time.Second)
		srvmanager.Log(*m.ErrorLogFile).Infof("%s ---- \n", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("%s\n", "log test ----")
		c.String(http.StatusOK, "Welcome !")
	})

	srv := &http.Server{
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	m.Run(srv)
}
```
