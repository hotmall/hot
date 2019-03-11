package main

import (
	"fmt"
	"os"
)

func main() {
	dirs := []string{
		"api",
		"code",
		"dist",
		"runtime",
		"runtime/bin",
		"runtime/etc/conf",
		"runtime/namedsql/mysql",
		"runtime/root",
		"runtime/var/log",
	}

	for _, dir := range dirs {
		fmt.Println(dir)
		os.MkdirAll(dir, os.ModeDir)
	}
}
