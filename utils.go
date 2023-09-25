/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"github.com/codeallergy/glue"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"strings"
)

/**
Checks that application started with verbose flag in command line.
 */

func IsVerbose(parent glue.Context) (verbose bool) {

	list := parent.Bean(ApplicationFlagsClass, glue.DefaultLevel)
	if len(list) > 0 {
		flags := list[0].Object().(ApplicationFlags)
		if flags.Verbose() {
			verbose = true
		}
	}

	return
}

/**
Filters child context list by role.

Uses direct match or suffix match.
Filter child context by suffix '-{role}' or '{role}'.
 */

func FilterChildrenByRole(parent glue.Context, role string) []glue.ChildContext {

	suffix :=  "-" + role

	var children []glue.ChildContext

	// filter child context by suffix '-{role}' or '{role}'
	for _, child := range parent.Children() {
		childRole := child.Role()
		if childRole == role || strings.HasSuffix(childRole, suffix) {
			children = append(children, child)
		}
	}

	return children
}

/**
Connects to the particular rpc server by using client role name and executes the callback function.
 */

func DoWithClientConn(parent glue.Context, role string, cb func(*grpc.ClientConn) error) error {

	verbose := IsVerbose(parent)
	if verbose {
		glue.Verbose(log.Default())
	}

	children := FilterChildrenByRole(parent, role)
	if len(children) != 1 {
		return errors.Errorf("application context should have only one client child context, but found '%d' for the role '%s'", len(children), role)
	}

	ctx, err := children[0].Object()
	if err != nil {
		return errors.Errorf("child '%s' context error, %v", role, err)
	}

	list := ctx.Bean(GrpcClientConnClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("client context should have one *grpc.ClientConn instance, but found '%d'", len(list))
	}
	bean := list[0]

	if client, ok := bean.Object().(*grpc.ClientConn); ok {
		return cb(client)
	} else {
		return errors.Errorf("invalid object '%v' found instead of *grpc.ClientConn in client context", bean.Class())
	}

}

/**
Connects to the control rpc server and executes callback function.
 */

func DoWithControlClient(parent glue.Context, cb func(ControlClient) error) error {

	return DoWithClient(parent, ControlClientRole, ControlClientClass, func(instance interface{}) error {
		if client, ok := instance.(ControlClient); ok {
			return cb(client)
		} else {
			return errors.Errorf("invalid object '%v' found instead of sprint.ControlClient in client context", reflect.TypeOf(instance).String())
		}
	})

}

/**
Connects to the rpc server and executes callback function.
*/

func DoWithClient(parent glue.Context, role string, clientType reflect.Type, cb func(clientInstance interface{}) error) error {

	verbose := IsVerbose(parent)
	if verbose {
		glue.Verbose(log.Default())
	}

	children := FilterChildrenByRole(parent, role)

	if len(children) != 1 {
		return errors.Errorf("application context should have only one child context for role '%s', but found '%d''", role, len(children), ControlClientRole)
	}

	ctx, err := children[0].Object()
	if err != nil {
		return errors.Errorf("child '%s' context error, %v", role, err)
	}

	list := ctx.Bean(clientType, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("client context should have one '%s' inference, but found '%d'", clientType.String(), len(list))
	}

	return cb(list[0])

}
