/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"crypto/tls"
	"github.com/codeallergy/glue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"net"
	"net/http"
	"reflect"
)


var (
	TlsConfigClass = reflect.TypeOf((*tls.Config)(nil)) // *tls.Config

	GrpcServerClass    = reflect.TypeOf((*grpc.Server)(nil))   // *grpc.Server
	HealthCheckerClass = reflect.TypeOf((*health.Server)(nil)) // *health.Server
	HttpServerClass    = reflect.TypeOf((*http.Server)(nil))   // *http.Server
)

var ServerClass = reflect.TypeOf((*Server)(nil)).Elem()

type EmptyAddrType struct {
}

func (t EmptyAddrType) Network() string {
	return ""
}

func (t EmptyAddrType) String() string {
	return ""
}

var EmptyAddr = EmptyAddrType{}

type Server interface {
	glue.InitializingBean
	glue.DisposableBean

	/**
	Bind server to the port.
	We separated it from the Serve, because we want to start application even if some servers were not able to bind.
	 */

	Bind() error

	/**
	Checks if server alive.
	 */

	Alive() bool

	/**
	Gets the actual listen address that could be different from bind address.
	The good example is if you bing to ip:0 it would have random port assigned to the socket.

	For non active server return EmptyAddr
	 */

	ListenAddress() net.Addr

	/**
	Runs actual server. The error code is the server exit code.
	We automatically filtering the 'closed' socket error codes, because they does not bring something valuable.
	 */

	Serve() error

	/**
	Shutdown server by the request.
	 */

	Shutdown() error

	/**
	ShutdownCh returns a channel that can be selected to wait
	for the server to perform a shutdown.
	 */

	ShutdownCh() <-chan struct{}

}

/**
Router interface for routing the HTTP request to specific pattern.
 */

var RouterClass = reflect.TypeOf((*Router)(nil)).Elem()

type Router interface {
	http.Handler

	/**
	Returns the url pattern used to serve the page.
	 */

	Pattern() string
}


