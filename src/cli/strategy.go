package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shurcooL/graphql"
)

// ExtensionContext holds execution context for extension operations
type ExtensionContext struct {
	Dir           string
	ManifestFile  string
	GraphqlClient *graphql.Client
}

// ExtensionLifecycle defines the interface for extension lifecycle operations
type ExtensionLifecycle interface {
	Build() (bool, error)
	Pack() (string, error)
	Publish() (bool, error)
	AddToRegistry(context *ExtensionContext) (bool, error)
}

// AbstractExtension provides common functionality for extension implementations
type AbstractExtension struct {
	ExtensionLifecycle

	Manifest *Manifest
	Context  *ExtensionContext
}

// NPMExtension implements lifecycle for NPM-based extensions
type NPMExtension struct {
	AbstractExtension
}

// RAWExtension implements lifecycle for simple non-NPM extensions
type RAWExtension struct {
	AbstractExtension
}

// Common Lifecycle commands

// Build executes the build process for an extension
func (ae AbstractExtension) Build() (bool, error) {
	return runCommand(ae.Context.Dir, ae.Manifest.Scripts.Build, "")
}

// AddToRegistry adds the extension to the CDEBase registry
func (ae AbstractExtension) AddToRegistry(ctx *ExtensionContext) (bool, error) {
	manifest := ae.Manifest.String()

	data, err := os.ReadFile(filepath.Join(ctx.Dir, ctx.ManifestFile))
	if err == nil {
		manifest = string(data)
	}

	fmt.Println("Preparing to add extension to registry...")
	fmt.Printf("Using endpoint: %s\n", config.Endpoint)
	fmt.Printf("Extension ID: %s\n", ae.Manifest.ExtensionID)

	mutation, variables := NewPublishExtensionMutation(PublishExtensionVariables{
		force:       true,
		manifest:    manifest,
		name:        ae.Manifest.Name,
		bundle:      ae.Manifest.Bundle,
		version:     ae.Manifest.Version,
		extensionID: ae.Manifest.ExtensionID,
	})

	// Set a context with timeout
	ctx2, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	netErr := ae.Context.GraphqlClient.Mutate(ctx2, &mutation, variables)
	if netErr != nil {
		fmt.Printf("GraphQL error: %v\n", netErr)
		return false, netErr
	}

	fmt.Printf("Extension %s successfully published\n", ae.Manifest.ExtensionID)
	return true, nil
}

// Pack prepares the extension for publishing
func (ae AbstractExtension) Pack() (string, error) {
	fmt.Println("Packing extension and preparing to publish...")

	if err := ae.Manifest.ReadAssets(ae.Context.Dir); err != nil {
		return "", err
	}

	if err := ae.Manifest.ReadBundle(ae.Context.Dir); err != nil {
		return "", err
	}

	return "", nil
}

// Implementation-specific methods

// Publish for RAW extensions
func (re RAWExtension) Publish() (bool, error) {
	fmt.Println("Publishing simple extension...")
	return true, nil
}

// Publish for NPM extensions
func (ne NPMExtension) Publish() (bool, error) {
	var cmd = fmt.Sprintf("npm publish --registry=%s", config.Registry)

	if ne.Manifest.Scripts.Publish != "" {
		cmd = ne.Manifest.Scripts.Publish
	}

	return runCommand(ne.Context.Dir, cmd, fmt.Sprintf("REGISTRY=%s", config.Registry))
}
