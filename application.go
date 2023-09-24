/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"context"
	"flag"
	"github.com/codeallergy/glue"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"reflect"
)

var (
	ZapLogClass = reflect.TypeOf((*zap.Logger)(nil)) // *zap.Logger
	HCLogClass = reflect.TypeOf((*hclog.Logger)(nil)).Elem() // hclog.Logger interface
	LumberjackClass = reflect.TypeOf((*lumberjack.Logger)(nil)) // *lumberjack.Logger

	FileSystemClass = reflect.TypeOf((*http.FileSystem)(nil)).Elem() // http.FileSystem
	FlagSetClass    = reflect.TypeOf((*flag.FlagSet)(nil))           // *flag.FlagSet
)

/**
Generic component class that has a name
*/
var ComponentClass = reflect.TypeOf((*Component)(nil)).Elem()

type Component interface {
	glue.NamedBean

	/**
	Gets status with name=value key pair.
	Server responds status request with stats ordered by key.
	 */

	GetStats(cb func(name, value string) bool) error
}


var ApplicationClass = reflect.TypeOf((*Application)(nil)).Elem()

/**
	Application is the base entry point class for golang application.
 */
type Application interface {
	context.Context
	glue.InitializingBean
	glue.NamedBean
	Component

	/**
	Add beans to application context
	 */
	AppendBeans(beans ...interface{})

	/**
	Gets application name, represents lower case normalized executable name
	 */
	Name() string

	/**
	Gets application version at the time of compilation
	*/
	Version() string

	/**
	Gets application build at the time of compilation
	*/
	Build() string

	/**
	Gets application runtime profile, could be: dev, qa, prod and etc.
	 */
	Profile() string

	/**
	Checks if application running in dev mode
	 */
	IsDev() bool

	/**
	Gets application binary name, used on startup, could be different with application name
	 */
	Executable() string

	/**
	Gets home directory of the application, usually parent directory of binary folder where executable is running, not current directory
	 */
	ApplicationDir() string

	/**
	Run application with command line arguments
	 */
	Run(args []string) error

	/**
	Indicator if application is active and not in shutting down mode
	 */
	Active() bool

	/**
	Sets the flag that application is in shutting down mode then notify all go routines by ShutdownChannel then notify signal channel with interrupt signal

	Additionally sets the flag that application is going to be restarted after shutdown
	 */
	Shutdown(restart bool)

	/**
	Indicator if application needs to be restarted by autoupdate or remote command after shutdown
	*/
	Restarting() bool

}

var SystemEnvironmentPropertyResolverClass = reflect.TypeOf((*SystemEnvironmentPropertyResolver)(nil)).Elem()

type SystemEnvironmentPropertyResolver interface {

	/**
	Prompt property from console
	 */

	PromptProperty(key string) (string, bool)

	/**
	Return all properties prompt early with values or without.
	 */

	Environ(withValues bool) []string

}

var ApplicationFlagsClass = reflect.TypeOf((*ApplicationFlags)(nil)).Elem()

type ApplicationFlags interface {
	glue.PropertyResolver

	/**
	Flag indicates that process is running in background mode.
	 */

	Daemon() bool

	/**
	Flag indicates that additional logging should post in console on startup,
	For example the context creation and beans initialization steps.
	 */

	Verbose() bool

	/**
	Returns the sequence number of the running node.
	*/

	Node() int

	/**
	Returns the override properties from console line
	 */

	Properties() map[string]string

}

var FlagSetRegistrarClass = reflect.TypeOf((*FlagSetRegistrar)(nil)).Elem()

type FlagSetRegistrar interface {

	/**
	Registers listener for FlagSet object in order to add custom flags.
	 */
	RegisterFlags(fs *flag.FlagSet)

	/**
	Adds additional args to restart the server.
	 */

	RegisterServerArgs(args []string) []string
}

