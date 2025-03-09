package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
	version  = "0.1.0"
	commands []*cli.Command
	flags    = ApplicationFlags{}
)

func main() {
	app := &cli.App{
		Version:     version,
		Description: "CLI Tool for CDEBase extensions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Path to config file",
				Destination: &flags.ConfigPath,
			},
			&cli.StringFlag{
				Name:        "registry",
				Usage:       "NPM Registry URL",
				Destination: &flags.Registry,
			},
			&cli.StringFlag{
				Name:        "access_token",
				Destination: &flags.AccessToken,
				Usage:       "Access Token for CDEBase Account",
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Destination: &flags.Endpoint,
				Usage:       "URL to CDEBase graphql server",
			},
			&cli.StringFlag{
				Name:        "manifest",
				Destination: &flags.Manifest,
				Usage:       "Manifest Filepath",
			},
			&cli.StringFlag{
				Name:        "dir",
				Destination: &flags.Dir,
				Usage:       "Extension directory",
			},
		},
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
