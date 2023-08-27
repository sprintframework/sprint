/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"github.com/codeallergy/glue"
	"google.golang.org/grpc"
	"io"
	"reflect"
)

var GrpcClientConnClass = reflect.TypeOf((*grpc.ClientConn)(nil))   // *grpc.ClientConn
var ClientScannerClass = reflect.TypeOf((*ClientScanner)(nil)).Elem()

type ClientScanner interface {

	ScannerName() string

	ClientBeans() []interface{}

}

var CommandClass = reflect.TypeOf((*Command)(nil)).Elem()

type Command interface {
	glue.NamedBean

	Run(args []string) error

	Desc() string
}

var ControlClientClass = reflect.TypeOf((*ControlClient)(nil)).Elem()

type ControlClient interface {
	glue.InitializingBean
	glue.DisposableBean

	Status() (string, error)

	Shutdown(restart bool) (string, error)

	ConfigCommand(command string, args []string) (string, error)

	CertificateCommand(command string, args []string) (string, error)

	JobCommand(command string, args []string) (string, error)

	StorageCommand(command string, args []string) (string, error)

	StorageConsole(writer io.StringWriter, errWriter io.StringWriter) error
}
