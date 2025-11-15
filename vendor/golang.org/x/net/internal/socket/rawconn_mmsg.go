// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package socket

import (
	"net"
)

func (c *Conn) recvMsgs(ms []Message, flags int) (int, error) {
	for i := range ms {
		ms[i].raceWrite()
	}
	tmps := defaultMmsgTmpsPool.Get()
	defer defaultMmsgTmpsPool.Put(tmps)
	var parseFn func([]byte, string) (net.Addr, error)
	if c.network != "tcp" {
		parseFn = parseInetAddr
	}
	hs := tmps.packer.pack(ms, parseFn, nil)
	n, err := tmps.syscaller.recvmmsg(c.c, hs, flags)
	if err != nil {
		return n, err
	}
	if err := hs[:n].unpack(ms[:n], parseFn, c.network); err != nil {
		return n, err
	}
	return n, nil
}

func (c *Conn) sendMsgs(ms []Message, flags int) (int, error) {
	for i := range ms {
		ms[i].raceRead()
	}
	tmps := defaultMmsgTmpsPool.Get()
	defer defaultMmsgTmpsPool.Put(tmps)
	var marshalFn func(net.Addr, []byte) int
	if c.network != "tcp" {
		marshalFn = marshalInetAddr
	}
	hs := tmps.packer.pack(ms, nil, marshalFn)
	n, err := tmps.syscaller.sendmmsg(c.c, hs, flags)
	if err != nil {
		return n, err
	}
	if err := hs[:n].unpack(ms[:n], nil, ""); err != nil {
		return n, err
	}
	return n, nil
}
