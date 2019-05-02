package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

type ApplicationFlags struct {
	Dir         string
	Manifest    string
	Registry    string
	Endpoint    string
	ConfigPath  string
	AccessToken string
}

var (
	config   *Config
	version  = "0.0.1"
	commands []cli.Command
	flags    = ApplicationFlags{}
)

func main() {
	var err error

	cli.BashCompletionFlag = cli.BoolFlag{
		Hidden: true,
		Name:   "compgen",
	}

	app := cli.NewApp()

	app.Version = version
	app.EnableBashCompletion = true
	app.Description = "CLI Tool for CDEBase extensions"

	if err != nil {
		log.Fatal(err)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "Path to config file",
			Destination: &flags.ConfigPath,
		},
		cli.StringFlag{
			Name:        "registry",
			Usage:       "NPM Registry URL",
			Destination: &flags.Registry,
		},
		cli.StringFlag{
			Name:        "access_token",
			Destination: &flags.AccessToken,
			Usage:       "Access Token for CDEBase Account",
		},
		cli.StringFlag{
			Name:        "endpoint",
			Destination: &flags.ConfigPath,
			Usage:       "URL to CDEBase graphql server",
		},
		cli.StringFlag{
			Name:        "manifest",
			Destination: &flags.Manifest,
			Usage:       "Manifest Filepath",
		},
		cli.StringFlag{
			Name:        "dir",
			Destination: &flags.Dir,
			Usage:       "Extension directory",
		},
	}

	// Add Commands
	app.Commands = commands

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
