package main

import (
	"fmt"
	"github.com/shurcooL/graphql"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

type Extension struct {
	dir      string
	entry    string
	bundle   string
	manifest *Manifest
}

func init() {
	commands = append(commands, cli.Command{
		Name:        "extension",
		Aliases:     []string{"e"},
		Description: "Extensions management",
		Subcommands: []cli.Command{
			{
				Name:        "publish",
				Aliases:     []string{"p"},
				Description: "Publish Extension to CDEBase Repository",
				Action: func(c *cli.Context) error {
					var err error
					//var manifest *Manifest
					var strategy ExtensionLifecycle

					pwd, _ := os.Getwd()
					fmt.Println("Publishing extension...")

					config, _ := loadConfig(flags.ConfigPath, &flags)
					client := graphql.NewClient(config.Endpoint, nil)

					context := ExtensionContext{
						GraphqlClient: client,
						ManifestFile:  flags.Manifest,
						Dir:           filepath.Join(pwd, flags.Dir),
					}

					_, strategy, err = ReadManifest(&context)

					_, err = strategy.Build()
					_, err = strategy.Pack()
					_, err = strategy.Publish()
					_, err = strategy.AddToRegistry(&context)

					if err != nil {
						log.Fatal(err)
					}

					return err
				},
			},
		},
	})
}
