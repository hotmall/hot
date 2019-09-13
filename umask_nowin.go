//+build !windows

package main

import (
	"syscall"
)

func umask(mask int)(oldmask int) {
	return syscall.Umask(mask)
}

