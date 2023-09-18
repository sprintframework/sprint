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

	/**
		Should return unique scanner name to debug and logging needs, the source of the bean.
	*/
	ScannerName() string

	/**
		List of all client beans.
	 */
	ClientBeans() []interface{}

}

var CommandClass = reflect.TypeOf((*Command)(nil)).Elem()

type Command interface {
	glue.NamedBean

    /**
	  Help should return long-form help text that includes the command-line usage,
	  a brief few sentences explaining the function of the command,
	  and the complete list of flags the command accepts.
     */
	Help() string

	/**
	  Run should run the actual command with the given CLI instance and command-line arguments.
	  It should return error or nil on success.
	*/
	Run(args []string) error

	/**
	  Synopsis should return a one-line, short synopsis of the command, should be less than 50 characters ideally.
	 */
	Synopsis() string
}

var ControlClientClass = reflect.TypeOf((*ControlClient)(nil)).Elem()

type ControlClient interface {
	glue.InitializingBean
	glue.DisposableBean

	/**
		Returned status of the node.
	 */
	Status() (string, error)

	/**
		Sends shutdown command with restart or not.
	 */
	Shutdown(restart bool) (string, error)

	/**
		Sends config command with options to set or get config entry.
	 */
	ConfigCommand(command string, args []string) (string, error)

	/**
		Sends certificate generation command with arguments.
	 */
	CertificateCommand(command string, args []string) (string, error)

	/**
		Sends job command with arguments.
	 */
	JobCommand(command string, args []string) (string, error)

	/**
		Sends storage command with arguments.
	 */
	StorageCommand(command string, args []string) (string, error)

	/**
		Initialize bi-directional console with defined writer and error writer streams.
	 */
	StorageConsole(writer io.StringWriter, errWriter io.StringWriter) error
}
