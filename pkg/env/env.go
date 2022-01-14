package env

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	// DefaultBhojpurDataRoot is the default value for BHOJPUR_ROOT environment variable
	DefaultBhojpurDataRoot = "/bhojpur"
	// DefaultBhojpurRoot is only required for hooks
	DefaultBhojpurRoot = "/usr/local/bhojpur"
)

// BhojpurRoot returns $BHOJPUR_ROOT or tries to guess its value if it's not set.
// It is a root for 'Bhojpur.NET' platform distribution, which contains bin/platform
// for instance.
func BhojpurRoot() (root string, err error) {
	if root = os.Getenv("BHOJPUR_ROOT"); root != "" {
		return root, nil
	}
	command, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	dir := path.Dir(command)

	if strings.HasSuffix(dir, "/bin") {
		return path.Dir(dir), nil
	}
	return DefaultBhojpurRoot, nil
}

// BhojpurDataRoot returns $BHOJPUR_DATAROOT or the default if $BHOJPUR_DATAROOT is not
// set. BhojpurDataRoot does not check if the directory exists and is
// writable.
func BhojpurDataRoot() string {
	if dataRoot := os.Getenv("BHOJPUR_DATAROOT"); dataRoot != "" {
		return dataRoot
	}

	return DefaultBhojpurDataRoot
}

// BhojpurMysqlRoot returns the root for the MySQL database distribution,
// which contains bin/mysql CLI for instance.
// If it is not set, look for mysqld in the path.
func BhojpurMysqlRoot() (string, error) {
	// if the environment variable is set, use that
	if root := os.Getenv("BHOJPUR_MYSQL_ROOT"); root != "" {
		return root, nil
	}

	// otherwise let's look for mysqld in the PATH.
	// ensure that /usr/sbin is included, as it might not be by default
	// This is the default location for mysqld from packages.
	newPath := fmt.Sprintf("/usr/sbin:%s", os.Getenv("PATH"))
	os.Setenv("PATH", newPath)
	path, err := exec.LookPath("mysqld")
	if err != nil {
		return "", errors.New("BHOJPUR_MYSQL_ROOT is not set and no mysqld could be found in your PATH")
	}
	path = filepath.Dir(filepath.Dir(path)) // strip mysqld, and the sbin
	return path, nil
}

// BhojpurMysqlBaseDir returns the MySQL base directory, which
// contains the fill_help_tables.sql script for instance
func BhojpurMysqlBaseDir() (string, error) {
	// if the environment variable is set, use that
	if root := os.Getenv("BHOJPUR_MYSQL_BASEDIR"); root != "" {
		return root, nil
	}

	// otherwise let's use BhojpurMysqlRoot
	root, err := BhojpurMysqlRoot()
	if err != nil {
		return "", errors.New("BHOJPUR_MYSQL_BASEDIR is not set. Please set $BHOJPUR_MYSQL_BASEDIR")
	}
	return root, nil
}
