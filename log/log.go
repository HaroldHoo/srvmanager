/**
* Copyright 2018 harold. All rights reserved.
* Filename: ./srvmanager/log.go
* Author: harold
* Mail: mail@yaolong.me
* Date: 2018-03-19
*/

package srvmanager

import(
	"log"
	"os"
	"fmt"
	"sync"
	"runtime"
)

const(
	P_FATAL		= "[FATAL] "
	P_ERROR		= "[ERROR] "
	P_WARNING	= "[WARNING] "
	P_INFO		= "[INFO] "
	P_DEBUG		= "[DEBUG] "
)

type logS struct{
	logger		*log.Logger
	file		*os.File
	filename	*string
}

var map_logs map[string]*logS
var mu sync.Mutex

func Logger(filename string, prefix string, flag int) (ret *log.Logger) {
	mu.Lock()
	defer mu.Unlock()

	if map_logs == nil{
		map_logs = make(map[string]*logS)
	}

	key := fmt.Sprintf("%+v|%+v|%+v", filename, prefix, flag)
	if map_logs[key] != nil {
		fS, _ := map_logs[key].file.Stat()
		if fS != nil{
			ret = map_logs[key].logger
			return
		}
	}

	file := new(os.File)
	err := new(error)
	if filename == "1" {
		file = os.Stdout
	} else if filename == "2"{
		file = os.Stderr
	} else {
		file, *err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if *err != nil{
			file.Close()
			Stderrf(2, P_FATAL, "%s, open log file faild.\n", *err)
			ret = nil
			return
		}
	}

	map_logs[key] = &logS{
		logger: log.New(file, prefix, flag),
		file: file,
		filename: &filename,
	}
	ret = map_logs[key].logger

	return
}

func Stdoutf(calldepth int, prefix string, format string, v ...interface{}){
	l := Logger("1", prefix, log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	if l != nil{
		format = fmt.Sprintf("pid:%d; %s", os.Getpid(), format)
		l.Output(calldepth, fmt.Sprintf(format, v...))
	}
}

func Stderrf(calldepth int, prefix string, format string, v ...interface{}){
	l := Logger("2", prefix, log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	if l != nil{
		format = fmt.Sprintf("pid:%d; %s", os.Getpid(), format)
		l.Output(calldepth, fmt.Sprintf(format, v...))
	}
}

func Writef(calldepth int, filename string, prefix string, format string, v ...interface{}){
	l := Logger(filename, prefix, log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	if l != nil{
		format = fmt.Sprintf("pid:%d; %s", os.Getpid(), format)
		l.Output(calldepth, fmt.Sprintf(format, v...))
	}
}

func GC(){
	for k,v := range map_logs{
		if v != nil {
			fs,_ := v.file.Stat()
			if fs == nil {
				v.file.Close()
				delete(map_logs, k)
			}
		}
	}
}

type Log struct{
	Filename	string
}

func (l *Log) Fatalf(format string, v ...interface{}){
	if l.Filename != "" {
		Writef(3, l.Filename, P_FATAL, format, v...)
		Writef(3, l.Filename, P_FATAL, "exit status 1\n")
	}

	os.Exit(1)
}

func (l *Log) Errorf(format string, v ...interface{}){
	if l.Filename != "" {
		Writef(3, l.Filename, P_ERROR, format, v...)
	}
}

func (l *Log) Warningf(format string, v ...interface{}){
	if l.Filename != "" {
		Writef(3, l.Filename, P_WARNING, format, v...)
	}
}

func (l *Log) Infof(format string, v ...interface{}){
	if l.Filename != "" {
		Writef(3, l.Filename, P_INFO, format, v...)
	}
}

func (l *Log) Debugf(format string, v ...interface{}){
	if l.Filename != "" {
		Writef(3, l.Filename, P_DEBUG, format, v...)
	}
}

func (l *Log) Stackf(all bool, prefix string, format string, v ...interface{}){
	format = fmt.Sprintf("%s; %s", format, string(Stack(all)))
	if l.Filename != "" {
		Writef(3, l.Filename, prefix, format, v...)
	}
}

func (l *Log) Stderrf(prefix string, format string, v ...interface{}){
	Stderrf(3, prefix, format, v...)
}

func (l *Log) Stdoutf(prefix string, format string, v ...interface{}){
	Stdoutf(3, prefix, format, v...)
}

func Stack(all bool) []byte {
	buf := make([]byte, 2048)
	for {
		n := runtime.Stack(buf, all)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

