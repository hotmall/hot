package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//log "github.com/Sirupsen/logrus"

//ClientCommand is executed to generate client from a RAML specification
type ClientCommand struct {
	Language    string
	Kind        string
	PackageName string
	Module      string
}

//Execute generates a client from a RAML specification
func (command *ClientCommand) Execute() error {
	//log.Debug("Generating a rest client for ", command.Language)

	var content []string
	content = append(content, fmt.Sprintf("package %s\n", command.PackageName))
	filepath.Walk("api", func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() && info.Name() == "types" {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, ".raml") {
			item := fmt.Sprintf("//go:generate go-raml client --language %s --kind %s --package %s --ramlfile %s --import-path %s", command.Language, command.Kind, command.PackageName, strings.Replace(path, "\\", "/", -1), command.Module)
			content = append(content, item)
		}
		return nil
	})
	content = append(content, "\n")

	fmt.Println("generate.go")
	ioutil.WriteFile("generate.go", []byte(strings.Join(content, "\n")), 0660)

	fmt.Println("go.mod")
	gomod := fmt.Sprintf(mod, command.Module, getGoVersion())
	ioutil.WriteFile("go.mod", []byte(gomod), 0660)

	return nil
}
