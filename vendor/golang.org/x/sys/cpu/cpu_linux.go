// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !386 && !amd64 && !amd64p32 && !arm64
// +build !386,!amd64,!amd64p32,!arm64

package cpu

func init() {
	if err := readHWCAP(); err != nil {
		return
	}
	doinit()
	Initialized = true
}
