package main

import (
	"fmt"
	"github.com/shurcooL/graphql"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Extension struct {
	dir      string
	entry    string
	bundle   string
	manifest *Manifest
}

func init() {
	commands = append(commands, &cli.Command{
		Name:        "extension",
		Aliases:     []string{"e"},
		Description: "Extensions management",
		Subcommands: []*cli.Command{
			{
				Name:        "publish",
				Aliases:     []string{"p"},
				Description: "Publish Extension to CDEBase Repository",
				Action: func(c *cli.Context) error {
					var err error
					var strategy ExtensionLifecycle

					pwd, _ := os.Getwd()
					fmt.Println("Publishing extension...")

					config, err := loadConfig(flags.ConfigPath, &flags)
					if err != nil {
						fmt.Printf("Error loading configuration: %v\n", err)
						return err
					}

					if config.Endpoint == "" {
						fmt.Println("Error: No endpoint specified in config or via flags")
						return fmt.Errorf("missing GraphQL endpoint in configuration")
					}

					// Create HTTP client with timeout
					httpClient := &http.Client{
						Timeout: 30 * time.Second,
					}
					
					// Create GraphQL client with the HTTP client
					client := graphql.NewClient(config.Endpoint, httpClient)

					context := ExtensionContext{
						GraphqlClient: client,
						ManifestFile:  flags.Manifest,
						Dir:           filepath.Join(pwd, flags.Dir),
						Config:        config,
					}

					fmt.Println("Reading manifest file:", context.ManifestFile)
					_, strategy, err = ReadManifest(&context)
					if err != nil {
						fmt.Printf("Error reading manifest: %v\n", err)
						return err
					}

					fmt.Println("Building extension...")
					if _, err = strategy.Build(); err != nil {
						fmt.Printf("Build failed: %v\n", err)
						return err
					}
					
					fmt.Println("Packing extension...")
					if _, err = strategy.Pack(); err != nil {
						fmt.Printf("Pack failed: %v\n", err)
						return err
					}
					
					fmt.Println("Publishing extension...")
					if _, err = strategy.Publish(); err != nil {
						fmt.Printf("Publish failed: %v\n", err)
						return err
					}
					
					fmt.Println("Adding to registry...")
					if _, err = strategy.AddToRegistry(&context); err != nil {
						fmt.Printf("Registry update failed: %v\n", err)
						return err
					}

					fmt.Println("Extension published successfully!")
					return nil
				},
			},
		},
	})
}
