/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package network

type Agent interface {
	Run()
	OnClose()
}
