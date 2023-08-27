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
)

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

func FindClientScanner(parent glue.Context, scannerName string) (ClientScanner, error) {

	list := parent.Bean(ClientScannerClass, glue.DefaultLevel)

	if scannerName == "" {
		if len(list) != 1 {
			return nil, errors.Errorf("application context should have one client scanner, but found '%d'", len(list))
		}

		bean := list[0]

		scanner, ok := bean.Object().(ClientScanner)
		if !ok {
			return nil, errors.Errorf("invalid object '%v' found instead of sprint.ClientScanner in application context", bean.Class())
		}

		return scanner, nil
	}

	for i, bean := range list {
		scanner, ok := bean.Object().(ClientScanner)
		if !ok {
			return nil, errors.Errorf("invalid object '%v' found on position %d instead of sprint.ClientScanner in application context", bean.Class(), i)
		}
		if scanner.ScannerName() == scannerName {
			return scanner, nil
		}
	}

	return nil, errors.Errorf("client scanner '%s' not found in application context", scannerName)

}

func DoWithClientConn(parent glue.Context, scannerName string, cb func(*grpc.ClientConn) error) error {

	verbose := IsVerbose(parent)
	scanner, err := FindClientScanner(parent, scannerName)
	if err != nil {
		return err
	}

	beans := scanner.ClientBeans()
	if verbose {
		verbose := glue.Verbose{ Log: log.Default() }
		beans = append([]interface{}{verbose}, beans...)
	}

	ctx, err := parent.Extend(beans...)
	if err != nil {
		return err
	}
	defer ctx.Close()

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

func DoWithControlClient(parent glue.Context, cb func(ControlClient) error) error {

	verbose := IsVerbose(parent)
	scanner, err := FindClientScanner(parent, "control")
	if err != nil {
		return err
	}

	beans := scanner.ClientBeans()
	if verbose {
		verbose := glue.Verbose{ Log: log.Default() }
		beans = append([]interface{}{verbose}, beans...)
	}

	ctx, err := parent.Extend(beans...)
	if err != nil {
		return err
	}
	defer ctx.Close()

	list := ctx.Bean(ControlClientClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("client context should have one sprint.ControlClient inference, but found '%d'", len(list))
	}
	bean := list[0]

	if client, ok := bean.Object().(ControlClient); ok {
		return cb(client)
	} else {
		return errors.Errorf("invalid object '%v' found instead of sprint.ControlClient in client context", bean.Class())
	}

}