package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//log "github.com/Sirupsen/logrus"
	"github.com/hotmall/hot/utils"
)

var gitIgnorePattern = `
dist/*.zip
dist/*.tgz
runtime/bin/%s
`

const hotPattern = `#!/bin/bash
APP_HOME=$(cd "$(dirname $0)";pwd)
cd $APP_HOME && hot server -l %s --kind %s --module %s
`

// ServerCommand is executed to generate a go server from a RAML specification
type ServerCommand struct {
	Language string // target language
	Kind     string
	Module   string
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
		"runtime/var/run",
	}

	mask := utils.Umask(0)
	defer utils.Umask(mask)

	for _, dir := range dirs {
		fmt.Println(dir)
		os.MkdirAll(dir, 0755)
	}

	if command.Language == "python" {
		fmt.Println("code/var/log")
		os.MkdirAll("code/var/log", 0755)

		fmt.Println("code/var/log/.gitignore")
		if !isFileExist("code/var/log/.gitignore") {
			os.WriteFile("code/var/log/.gitignore", []byte("*.out"), 0660)
		}
	}

	var content []string
	if command.Language == "go" {
		content = append(content, "package main\n")
	} else {
		content = append(content, "#!/bin/bash\n")
	}

	filepath.Walk("api", func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		fmt.Println("--", path)
		if info.IsDir() {
			if info.Name() == path {
				// 如果是 root 目录，返回 nil
				return nil
			}
			// 所有目录下的 raml 都忽略，比如 types, traits 等
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, ".raml") {
			if command.Language == "go" {
				content = append(content, fmt.Sprintf("//go:generate go-raml server --language %s --kind %s --ramlfile ../%s --no-apidocs --import-path %s", command.Language, command.Kind, strings.Replace(path, "\\", "/", -1), command.Module))
			} else {
				content = append(content, fmt.Sprintf("go-raml server --language %s --kind %s --ramlfile ../%s --no-apidocs --import-path %s", command.Language, command.Kind, strings.Replace(path, "\\", "/", -1), command.Module))
			}
		}
		return nil
	})
	// content = append(content, "\n")

	if command.Language == "python" {
		content = append(content, pythonInit)

		if command.Kind == "flask" {
			start := fmt.Sprintf(gunicornStart, command.Module)
			os.WriteFile("runtime/bin/start.sh", []byte(start), 0660)

			stop := fmt.Sprintf(gunicornStop, command.Module)
			os.WriteFile("runtime/bin/stop.sh", []byte(stop), 0660)
		}

		// create virtualenv setup.sh
		os.WriteFile("runtime/bin/setup.sh", []byte(venvSetup), 0660)

		gitIgnorePattern = pyGitIgnorePattern
	}

	if command.Language == "go" {
		if !isFileExist("runtime/bin/start.sh") {
			start := fmt.Sprintf(goStart, command.Module)
			os.WriteFile("runtime/bin/start.sh", []byte(start), 0660)
		}

		stop := fmt.Sprintf(goStop, command.Module)
		os.WriteFile("runtime/bin/stop.sh", []byte(stop), 0660)

		fmt.Println("code/generate.go")
		generate := strings.Join(content, "\n")
		generate += "\n"
		os.WriteFile("code/generate.go", []byte(generate), 0660)

		fmt.Println("code/go.mod")
		gomod := fmt.Sprintf(mod, command.Module, getGoVersion())
		if !isFileExist("code/go.mod") {
			os.WriteFile("code/go.mod", []byte(gomod), 0660)
		}
	} else {
		fmt.Println("code/generate.sh")
		os.WriteFile("code/generate.sh", []byte(strings.Join(content, "\n")), 0660)
	}

	fmt.Println("code/VERSION")
	if !isFileExist("code/VERSION") {
		os.WriteFile("code/VERSION", []byte("0.1.0"), 0660)
	}

	fmt.Println(".gitignore")
	exeName := filepath.Base(command.Module)
	ignore := fmt.Sprintf(gitIgnorePattern, exeName)
	if !isFileExist(".gitignore") {
		os.WriteFile(".gitignore", []byte(ignore), 0660)
	}

	fmt.Println("hot.sh")
	hot := fmt.Sprintf(hotPattern, command.Language, command.Kind, command.Module)
	os.WriteFile("hot.sh", []byte(hot), 0660)

	fmt.Println("runtime/var/log/.gitignore")
	if !isFileExist("runtime/var/log/.gitignore") {
		os.WriteFile("runtime/var/log/.gitignore", []byte("*.out\n*.log"), 0660)
	}

	fmt.Println("runtime/var/run/.gitignore")
	if !isFileExist("runtime/var/run/.gitignore") {
		os.WriteFile("runtime/var/run/.gitignore", []byte("*.pid"), 0660)
	}

	fmt.Println("dist/.gitignore")
	if !isFileExist("dist/.gitignore") {
		os.WriteFile("dist/.gitignore", []byte("*.tgz\n*.zip\n*.gz"), 0660)
	}

	return nil
}

// cek if a file exist
func isFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
