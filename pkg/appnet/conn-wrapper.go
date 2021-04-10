// Copyright 2020 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package appnet provides a simplified and functionally extended wrapper interface to the
scionproto/scion package snet.


Dispatcher and SCION daemon connections

During the hidden initialisation of this package, the dispatcher and sciond
connections are opened. The sciond connection determines the local IA.
The dispatcher and sciond sockets are assumed to be at default locations, but this can
be overridden using environment variables:

		SCION_DISPATCHER_SOCKET: /run/shm/dispatcher/default.sock
		SCION_DAEMON_ADDRESS: 127.0.0.1:30255

This is convenient for the normal use case of running the endhost stack for a
single SCION AS. When running multiple local ASes, e.g. during development, the
address of the sciond corresponding to the desired AS needs to be specified in
the SCION_DAEMON_ADDRESS environment variable.


Wildcard IP Addresses

snet does not currently support binding to wildcard addresses. This will hopefully be
added soon-ish, but in the meantime, this package emulates this functionality.
There is one restriction, that applies to hosts with multiple IP addresses in the AS:
the behaviour will be that of binding to one specific local IP address, which means that
the application will not be reachable using any of the other IP addresses.
Traffic sent will always appear to originate from this specific IP address,
even if that's not the correct route to a destination in the local AS.

This restriction will very likely not cause any issues, as a fairly contrived
network setup would be required. Also, sciond has a similar restriction (binds
to one specific IP address).
*/
package appnet

import (
	"C"
	"net"

	"github.com/scionproto/scion/go/lib/snet"
)
import (
	"errors"
	"time"
)

type PathWrapper struct {
    path snet.Path
}

type PathMetadataWrapper struct {
    meta *snet.PathMetadata
}

type PathListWrapper struct {
	paths []snet.Path
}

func (w PathListWrapper)Count() int {
	return len(w.paths)
}

func (w PathListWrapper)GetPathAt(index int) *PathWrapper {
	return &PathWrapper{ w.paths[index] }
}

func QueryPathsWrapped(addr *AddressWrapper) (*PathListWrapper, error) {
	udpAddr, ok := addr.addr.(*snet.UDPAddr)

	if !ok {
		return nil, errors.New("Can't look up paths for address that is not snet.UDPAddr")
	}

	spaths, err := QueryPaths(udpAddr.IA)
	if err != nil {
		return nil, err
	}
	
    return &PathListWrapper { spaths }, nil
}

// In bytes
func (m PathMetadataWrapper)GetMTU() int32 {
    return int32(m.meta.MTU);
} 

// In microseconds
func (m PathMetadataWrapper)GetLatencyAt(index int) int32 {
    return int32(time.Duration(m.meta.Latency[index])*time.Microsecond)
}

func (m PathWrapper)Length() int {
	return len(m.path.Path().Raw)
}

func (m PathWrapper)GetFingerprint() string {
	return string(snet.Fingerprint(m.path))
}

// In kbit/s
func (m PathMetadataWrapper)GetBandwidthAt(index int) int64 {
    return int64(m.meta.Bandwidth[index])
}

// Unix timestamp in s at UTC
func (m PathMetadataWrapper)GetExpiry() int64 {
    return m.meta.Expiry.UTC().Unix()
}

func (p PathWrapper)GetMetadata() *PathMetadataWrapper {
    return &PathMetadataWrapper { p.path.Metadata() }
}

// Wrapper around snet.conn to provide ObjC interoperability
type ConnWrapper struct {
	conn *snet.Conn
}

// Wrapper around net.Addr to provide ObjC interoperability
type AddressWrapper struct {
    addr net.Addr
}

// More interop stuff..
type ReadResult struct {
    BytesRead int
    Source *AddressWrapper
    Err error
}

func DialWrappedWithPath(address string, path *PathWrapper) (*ConnWrapper, error) {
	raddr, err := ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	if path != nil {
		SetPath(raddr, path.path)
	}
	c, e :=  DialAddr(raddr)
    return &ConnWrapper{c}, e
}

func DialWrapped(address string) (*ConnWrapper, error) {
    return DialWrappedWithPath(address, nil)
}

func ListenPortWrapped(port int) (*ConnWrapper, error) {
    c, e := ListenPort(uint16(port))
    return &ConnWrapper{c}, e
}

func (w ConnWrapper) LocalAddress() *AddressWrapper {
    return &AddressWrapper{w.conn.LocalScionAddr()};
}

func (w ConnWrapper) Read(buffer []byte) *ReadResult {
    n, a, e := w.conn.ReadFrom(buffer)
    
    return &ReadResult{n, &AddressWrapper{a}, e}
}

func (w ConnWrapper) Write(buffer []byte) (int, error) {
    return w.conn.Write(buffer)
}

func (w ConnWrapper) WriteTo(buffer []byte, address *AddressWrapper) (int, error) {
    return w.conn.WriteTo(buffer, address.addr)
}

func (w ConnWrapper) Close() {
    w.conn.Close()
}

func (w AddressWrapper) AsString() string {
    return w.addr.String()
}

func AddressWrapperFromString(str string) (*AddressWrapper, error) {
    raddr, err := ResolveUDPAddr(str)
	if err != nil {
		return nil, err
	}

    return &AddressWrapper{raddr}, nil
}
