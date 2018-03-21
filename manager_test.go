/**
 * Copyright 2018 harold. All rights reserved.
 * Filename: ./srvmanager/manager_test.go
 * Author: harold
 * Mail: mail@yaolong.me
 * Date: 2018-03-20
 */

package srvmanager

import (
	"fmt"
	"testing"
)

func TestGetCurrentExecname(t *testing.T) {
	fmt.Printf("\n----------------- TestGetCurrentExecname\n")
	var s1 = GetCurrentExecname()
	var s2 = GetCurrentExecname()
	var s3 = *currentExecname
	if s1 != s2 || s2 != s3 || s3 != s1 {
		t.Error(s1, s2, s3)
	}
	if s1 == "" {
		t.Error(s1, s2, s3)
	}
	fmt.Printf("%#v\n%#v\n%#v\n", s1, s2, s3)
}
