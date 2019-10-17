package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//log "github.com/Sirupsen/logrus"
	"github.com/hotmall/orange/utils"
)

// ServerCommand is executed to generate a go server from a RAML specification
type ServerCommand struct {
	Language string // target language
	Kind     string
}

// Execute generates a Go server from an RAML specification
func (command *ServerCommand) Execute() error {

	//log.Infof("Generating a %v server", command.Language)

	dirs := []string{
		"api",
		"code",
		"dist",
		"runtime",
		"runtime/bin",
		"runtime/etc/conf",
		"runtime/namedsql",
		"runtime/root",
		"runtime/var/log",
	}

	mask := utils.Umask(0)
	defer utils.Umask(mask)

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
			content = append(content, fmt.Sprintf("//go:generate go-raml server --language %s --kind %s --ramlfile ../%s --no-apidocs", command.Language, command.Kind, strings.Replace(path, "\\", "/", -1)))
		}
		return nil
	})
	content = append(content, "\n")

	ioutil.WriteFile("code/generate.go", []byte(strings.Join(content, "\n")), 0660)
	ioutil.WriteFile("code/VERSION", []byte("0.1.0"), 0660)

	return nil
}
