/**
 * Copyright 2018 harold. All rights reserved.
 * Filename: ./srvmanager/pidfile.go
 * Author: harold
 * Mail: mail@yaolong.me
 * Date: 2018-03-16
 */

package srvmanager

import (
	srv_log "github.com/HaroldHoo/srvmanager/log"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"errors"
)

func (m *Manager) getPidFromPidFile() (ret int, err error) {
	if m.PidFile == nil {
		panic("pidfile should not be empty.")
	}
	f, err := os.OpenFile(*m.PidFile, os.O_RDONLY, 0)
	defer f.Close()
	CheckErrAndExitToStderr(err)

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB)
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	CheckErrAndExitToStderr(err)

	b := make([]byte, 8)
	var n int; n, err = f.Read(b)
	if n == 0 {
		srv_log.Stderrf(3, srv_log.P_FATAL, "pidfile(%s) is empty; exit status 1\n", *m.PidFile)
		os.Exit(1)
	}
	CheckErrAndExitToStderr(err)

	ret, err = strconv.Atoi(strings.TrimRight(string(b), "\x00"))
	CheckErrAndExitToStderr(err)

	return
}

func (m *Manager) writePidFile() {
	if m.PidFile == nil {
		panic("pidfile should not be empty.")
	}

	f, err := os.Create(*m.PidFile)
	defer f.Close()
	CheckErrAndExitToStderr(err)

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	CheckErrAndExitToStderr(err)

	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%d", os.Getpid()))
	w.Flush()
}

func (m *Manager) removePidFile() (err error){
	err = errors.New("")
	if m.PidFile == nil {
		return
	}
	if fileIsExist(m.PidFile) {
		os.Remove(*m.PidFile)
	}
	return
}

func fileIsExist(filename *string) (ret bool) {
	ret = true
	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		ret = false
	}
	return
}

