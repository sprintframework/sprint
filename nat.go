/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"net"
	"reflect"
	"time"
)

var NatServiceClass = reflect.TypeOf((*NatService)(nil)).Elem()

type NatService interface {

	AllowMapping() bool

	AddMapping(protocol string, extport, intport int, name string, lifetime time.Duration) error

	DeleteMapping(protocol string, extport, intport int) error

	ExternalIP() (net.IP, error)

	ServiceName() string
}
