/**
 * Copyright 2018 harold. All rights reserved.
 * Filename: ./srvmanager/log_test.go
 * Author: harold
 * Mail: mail@yaolong.me
 * Date: 2018-03-19
 */

package srvmanager

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Printf("\n----------------- TestLog\n")
	Stdoutf(2, P_ERROR, "%s\n", "test1")
	fmt.Printf("map: %+v\n", map_logs)
	Stdoutf(2, P_ERROR, "%s\n", "test2")
	fmt.Printf("map: %+v\n", map_logs)
}

func TestFile(t *testing.T) {
	fmt.Printf("\n----------------- TestFile\n")
	Writef(2, "/tmp/testF", P_WARNING, "%s\n", "file test - ")
	fmt.Printf("map_logs: %+v\n", map_logs)

	Writef(2, "/tmp/testF", P_WARNING, "%s\n", "file test - ")

	fmt.Printf("map_logs: %+v\n", map_logs)
}

func TestGC(t *testing.T) {
	fmt.Printf("\n----------------- TestGC\n")

	map_logs = nil
	Writef(2, "/tmp/testG", P_WARNING, "%s\n", "file test - ")
	fmt.Printf("1 - map_logs: %+v\n", map_logs)

	for _, v := range map_logs {
		v.file.Close()
	}

	Writef(2, "/tmp/testG", P_WARNING, "%s\n", "file test - ")

	GC()
	fmt.Printf("2 - map_logs: %+v\n", map_logs)

	Writef(2, "/tmp/testG", P_WARNING, "%s\n", "file test - ")

	fmt.Printf("3 - map_logs: %+v\n", map_logs)
}

func TestErrorf(t *testing.T) {
	fmt.Printf("\n----------------- TestErrorf\n")
	l := &Log{Filename: "/tmp/testE"}
	l.Errorf("%s\n", "xxxxxxxxxxxx")
}

func TestFatalf(t *testing.T) {
	t.Skip()
	fmt.Printf("\n----------------- TestFatalf\n")
	l := &Log{Filename: "/tmp/testF"}
	l.Fatalf("%s\n", "xxxxxxxxxxxx")
}
