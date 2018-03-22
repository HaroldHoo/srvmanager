/**
* Copyright 2018 harold. All rights reserved.
* Filename: manager.go
* Author: harold
* Mail: mail@yaolong.me
* Date: 2018-03-15
 */

package srvmanager

import (
	srv_log "github.com/HaroldHoo/srvmanager/log"
	"log"
	"context"
	"flag"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"errors"
	"github.com/fvbock/endless"
)

type Manager struct {
	Srv           *http.Server
	PidFile       *string
	ErrorLogFile  *string
	AccessLogFile *string
	flag_signal   *string
	flag_port     *string
	flag_deamon   *bool
}

func New() *Manager {
	if flag.Parsed() {
		panic("flag must not parsed.")
	}

	m := &Manager{}

	pidfile := flag.String("pid", "/var/run/"+GetCurrentExecname()+".pid", "pid file")
	errorlogfile := flag.String("errorlog", "/var/log/"+GetCurrentExecname()+"_error.log", "log file")
	accesslogfile := flag.String("accesslog", "/var/log/"+GetCurrentExecname()+"_access.log", "log file")
	m.flag_deamon = flag.Bool("d", true, "Start as deamon.")
	m.flag_signal = flag.String("s", "start", "Send a signal to the process.  The argument signal can be one of: start stop reload restart,\nThe following table shows the corresponding system signals:\nstop	SIGTERM\nreload	SIGHUP\nrestart	SIGHUP\n")
	m.flag_port = flag.String("p", "8080", "Listen port")
	flag.Parse()

	if m.ErrorLogFile == nil {
		ferr := CheckFileCRW(*errorlogfile)
		CheckErrAndExitToStderr(ferr)
		m.ErrorLogFile = errorlogfile
	}
	if m.AccessLogFile == nil {
		ferr := CheckFileCRW(*accesslogfile)
		CheckErrAndExitToStderr(ferr)
		m.AccessLogFile = accesslogfile
	}
	if m.PidFile == nil {
		ferr := CheckFileCRW(*pidfile)
		CheckErrAndExitToStderr(ferr)
		m.PidFile = pidfile
	}

	return m
}

func (m *Manager) Run(server *http.Server) {
	if *m.flag_deamon == false {
		m.Srv = server
		m.Srv.Addr = ":" + *m.flag_port
		m.runServer()
	} else {
		switch *m.flag_signal {
		case "stop":
			pid,_ := m.getPidFromPidFile()
			syscall.Kill(pid, syscall.SIGINT)
		case "restart", "reload":
			pid,_ := m.getPidFromPidFile()
			err := syscall.Kill(pid, syscall.SIGHUP)
			if err != nil {
				m.startNewServer()
			}
		default:
			m.startNewServer()
		}
	}
}

func (m *Manager) startNewServer() {
	log := &srv_log.Log{Filename: *m.ErrorLogFile}

	path := GetCurrentFilename()
	argv := os.Args
	// log.Infof("exec path: %v\n", path)
	log.Infof("origin args: %s\n", argv)

	newArgs := make([]string, 0)
	newArgs = append(newArgs, os.Args[0])
	newArgs = append(newArgs, "-d=false")
	if m.flag_port != nil {
		newArgs = append(newArgs, "-p=" + *m.flag_port)
	}
	if m.PidFile != nil {
		newArgs = append(newArgs, "-pid=" + *m.PidFile)
	}
	if m.ErrorLogFile != nil {
		newArgs = append(newArgs, "-errorlog=" + *m.ErrorLogFile)
	}
	if m.AccessLogFile != nil {
		newArgs = append(newArgs, "-accesslog=" + *m.AccessLogFile)
	}
	log.Infof("deamon args: %s\n", newArgs)

	cmd := exec.Command(path)
	cmd.Args = newArgs
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Start: Failed to launch, error: %v\n", err)
	}
}

func (m *Manager) gracefulShutdown(second int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(second)*time.Second)
	defer cancel()
	if err := m.Srv.Shutdown(ctx); err != nil {
		log := &srv_log.Log{Filename: *m.ErrorLogFile}
		log.Fatalf("Server Shutdown:", err)
	}
}

func (m *Manager) runServer() {
	l := &srv_log.Log{Filename: *m.ErrorLogFile}
	l.Infof("Server's pid: %d\n", os.Getpid())
	m.writePidFile()

	// endless
	log.SetOutput(m.GetErrorLogWriter())
	log.SetPrefix("[LOG] ")
	endless.DefaultReadTimeOut = m.Srv.ReadTimeout
	endless.DefaultWriteTimeOut = m.Srv.WriteTimeout
	endless.DefaultMaxHeaderBytes = m.Srv.MaxHeaderBytes

	if err := endless.ListenAndServe(m.Srv.Addr, m.Srv.Handler); err != nil {
		if err != http.ErrServerClosed{
			m.removePidFile()
			l.Fatalf("%s\n", err)
		}else{
			l.Errorf("%s\n", err)
		}
	}
}

func (m *Manager) GetAccessLogWriter() (file *os.File){
	err := errors.New("")
	file, err = os.OpenFile(*m.AccessLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	if err != nil{
		file.Close()
		file = os.Stderr
	}
	return
}

func (m *Manager) GetErrorLogWriter() (file *os.File){
	err := errors.New("")
	file, err = os.OpenFile(*m.ErrorLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	if err != nil{
		file.Close()
		file = os.Stderr
	}
	return
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	PanicErr(err)
	return strings.Replace(dir, "\\", "/", -1)
}

func GetCurrentFilename() string {
	dir, err := filepath.Abs(os.Args[0])
	PanicErr(err)
	return strings.Replace(dir, "\\", "/", -1)
}

var currentExecname = new(string)

func GetCurrentExecname() (ret string) {
	if *currentExecname != "" {
		ret = *currentExecname
		return
	}
	str := GetCurrentFilename()
	*currentExecname = str[strings.LastIndex(str, "/")+1 : len(str)]
	ret = *currentExecname
	return
}

func CheckFileCRW(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	return err
}

func CheckErrAndExitToStderr(err error) {
	if err != nil {
		srv_log.Stderrf(3, srv_log.P_FATAL, "%s; exit status 1\n", err)
		os.Exit(1)
	}
}

func Log(filename string) (*srv_log.Log) {
	return &srv_log.Log{Filename: filename}
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

