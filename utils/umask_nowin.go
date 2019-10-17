//+build !windows

package utils

import (
	"syscall"
)

// Umask call syscall.Umask
func Umask(mask int) (oldmask int) {
	return syscall.Umask(mask)
}
