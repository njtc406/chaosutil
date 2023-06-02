/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package main

import (
	"os"
	"os/signal"
	"syscall"
)

var closeCh = make(chan os.Signal, 0)

func init() {
	signal.Notify(closeCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
}

func main() {
	serverByetcd()
}
