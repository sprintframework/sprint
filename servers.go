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

var ServerScannerClass = reflect.TypeOf((*ServerScanner)(nil)).Elem()

type ServerScanner interface {

	ServerBeans() []interface{}

}

var ServerClass = reflect.TypeOf((*Server)(nil)).Elem()

type Server interface {
	glue.InitializingBean
	glue.DisposableBean

	Bind() error

	Active() bool

	ListenAddress() net.Addr

	Serve() error

	Stop()
}

/**
Page interface for html pages
 */

var PageClass = reflect.TypeOf((*Page)(nil)).Elem()

type Page interface {
	http.Handler

	Pattern() string
}


