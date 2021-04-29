package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hotmall/hot/commands"
	"github.com/urfave/cli"
)

var (
	//ApplicationName is the name of the application
	ApplicationName = "Hot"
)

var (
	serverCommand = &commands.ServerCommand{}
	clientCommand = &commands.ClientCommand{}
)

func getModule() string {
	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		fmt.Println("Fail to get env GOPATH, err=", err)
		return ""
	}
	gopath = filepath.Join(gopath, "src")

	currDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Fail to get curr dir, err=", err)
		return ""
	}

	if !strings.HasPrefix(currDir, gopath) {
		return filepath.Base(currDir)
	}

	// set module
	module, err := filepath.Rel(gopath, currDir)
	if err != nil {
		panic("Failed to set module automatically:" + err.Error())
	}

	// re-join because otherwise windows will use `\`
	return path.Join(strings.Split(module, string(filepath.Separator))...)
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version: \t\t%v\nCommit Hash: \t\t%v\nBuild Date: \t\t%v\nGo Version: \t\t%v\n",
			commands.Version, commands.CommitHash, commands.BuildDate, commands.GoVersion)
	}

	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = commands.Version
	app.Usage = "Generate a generate.go file from the raml files in the api directory."

	module := getModule()
	defaultPackage := filepath.Base(module)

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Generate a server according to a RAML specification",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "language, l",
					Value:       "go",
					Usage:       "Language to construct a server for",
					Destination: &serverCommand.Language,
				},
				cli.StringFlag{
					Name:        "kind",
					Value:       "gorestful",
					Usage:       "Kind of server to generate (gorestful)",
					Destination: &serverCommand.Kind,
				},
				cli.StringFlag{
					Name:        "module",
					Value:       module,
					Usage:       "Module name for go mod",
					Destination: &serverCommand.Module,
				},
			},
			Action: func(c *cli.Context) {
				if err := serverCommand.Execute(); err != nil {
					//log.Error(err)
					fmt.Println(err)
				}
			},
		},
		{
			Name:  "client",
			Usage: "Create a client for a RAML specification",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "language, l",
					Value:       "go",
					Usage:       "Language to construct a client for",
					Destination: &clientCommand.Language,
				},

				cli.StringFlag{
					Name:        "kind",
					Value:       "requests",
					Usage:       "Kind of client to generate (requests,grequests)",
					Destination: &clientCommand.Kind,
				},
				cli.StringFlag{
					Name:        "package",
					Value:       defaultPackage,
					Usage:       "package name",
					Destination: &clientCommand.PackageName,
				},
				cli.StringFlag{
					Name:        "module",
					Value:       module,
					Usage:       "Module name for go mod",
					Destination: &clientCommand.Module,
				},
			},
			Action: func(c *cli.Context) {
				if err := clientCommand.Execute(); err != nil {
					//log.Error(err)
					fmt.Println(err)
				}
			},
		},
	}
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
