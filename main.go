package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	mask := umask(0)
	defer umask(mask)

	for _, dir := range dirs {
		fmt.Println(dir)
		os.MkdirAll(dir, 0755)
	}

	var content []string
	content = append(content, "package main\n")
	filepath.Walk("api", func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() && info.Name() == "types" {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, ".raml") {
			content = append(content, fmt.Sprintf("//go:generate go-raml server --kind gorestful --ramlfile ../%s --no-apidocs", path))
		}
		return nil
	})
	content = append(content, "\n")

	ioutil.WriteFile("code/generate.go", []byte(strings.Join(content, "\n")), 0660)
	ioutil.WriteFile("code/VERSION", []byte("0.1.0"), 0660)
}
