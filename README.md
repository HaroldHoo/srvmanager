# srvmanager
1. graceful shutdown/reload
2. manager process via pidfile
3. log something easier

## Examples
   1. Via https://github.com/gin-gonic/gin :
   > https://github.com/HaroldHoo/srvmanager/tree/master/examples/gin

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

