// Copyright 2018-2022 (c) The Go Service Components authors. All rights reserved. Issued under the Apache 2.0 License.

package network // import "github.com/karlmutch/go-service/pkg/network"

import (
	"net"
	"os"

	"github.com/go-stack/stack"
	"github.com/karlmutch/kv" // MIT License

	"github.com/karlmutch/go-fqdn"
)

// Functions related to networking needs for the server

// GetFreePort will find and return a port number that is found to be available
func GetFreePort(hint string) (port int, err kv.Error) {
	addr, errGo := net.ResolveTCPAddr("tcp", hint)
	if errGo != nil {
		return 0, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime())
	}

	l, errGo := net.ListenTCP("tcp", addr)
	if errGo != nil {
		return 0, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime())
	}

	port = l.Addr().(*net.TCPAddr).Port

	// Dont defer as the port will be quickly reused
	l.Close()

	return port, nil
}

// GetHostName returns a human readable host name that contains as much useful context
// as can be gathered
func GetHostName() (name string) {

	name = fqdn.Get()
	if len(name) != 0 && name != "unknown" {
		return name
	}

	name, _ = os.Hostname()

	if len(name) != 0 {
		return name
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return (ipnet.IP.String())
			}
		}
	}
	return "unknown"
}
