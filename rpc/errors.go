/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package rpc

import (
	"errors"
)

var (
	ErrServiceHasRegistered = errors.New("the service has already been registered")
)
