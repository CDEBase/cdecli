package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	ASSET_MAP       = "asset.map"
	ASSET_ICON      = "asset.icon"
	ASSET_README    = "asset.readme"
	ASSET_LICENSE   = "asset.license"
	ASSET_CHANGELOG = "asset.changelog"

	TYPE_SERVER  = "extension.server"
	TYPE_SIMPLE  = "extension.simple"
	TYPE_BROWSER = "extension.browser"
	TYPE_COMPLEX = "extension.complex"
)

var (
	EXTENSION_ASSETS = make(map[string][]string)
)

func init() {
	// Initialize extension assets map
	EXTENSION_ASSETS[ASSET_MAP] = []string{".map"}
	EXTENSION_ASSETS[ASSET_ICON] = []string{"icon.png", "icon.ico", "favicon.ico"}
	EXTENSION_ASSETS[ASSET_LICENSE] = []string{"LICENSE.md", "license.md", "LICENSE.txt", "LICENSE", "license.txt", "license", "License.md", "License.txt", "License"}
	EXTENSION_ASSETS[ASSET_README] = []string{"README.md", "readme.md", "README.txt", "README", "readme.md", "readme.txt", "readme", "Readme.md", "Readme.txt", "Readme"}
	EXTENSION_ASSETS[ASSET_CHANGELOG] = []string{"CHANGELOG.md", "changelog.md", "CHANGELOG.txt", "CHANGELOG", "changelog.txt", "changelog", "Changelog.md", "Changelog.txt", "Changelog"}
}

func ReadManifest(context *ExtensionContext) (*Manifest, ExtensionLifecycle, error) {
	manifest := Manifest{}
	var strategy ExtensionLifecycle

	data, err := os.ReadFile(filepath.Join(context.Dir, context.ManifestFile))
	if err != nil {
		return nil, nil, fmt.Errorf("%s\n\nRun this command in a directory with a %s file for an extension.\n\n", err, context.ManifestFile)
	}

	if err = json.Unmarshal(data, &manifest); err != nil {
		return nil, nil, err
	}

	if manifest.Extension.ExtensionType == TYPE_SIMPLE {
		strategy = RAWExtension{
			AbstractExtension{
				Context:  context,
				Manifest: &manifest,
			},
		}
	} else {
		strategy = NPMExtension{
			AbstractExtension{
				Context:  context,
				Manifest: &manifest,
			},
		}
	}

	return &manifest, strategy, nil
}

type NpmPackage struct {
	Icon                  string            `json:"icon"`
	Name                  string            `json:"name"`
	Main                  string            `json:"main"`
	Type                  string            `json:"type"`
	Badges                []string          `json:"badges"`
	Preview               bool              `json:"preview"`
	Version               string            `json:"version"`
	Keywords              []string          `json:"keywords"`
	License               string            `json:"license"`
	Homepage              string            `json:"homepage"`
	Publisher             string            `json:"publisher"`
	Categories            []string          `json:"categories"`
	Description           string            `json:"description"`
	ExtensionID           string            `json:"extensionID"`
	DisplayName           string            `json:"displayName"`
	Dependencies          map[string]string `json:"dependencies"`
	DevDependencies       map[string]string `json:"devDependencies"`
	ActivationEvents      []string          `json:"activationEvents"`
	ExtensionDependencies []string          `json:"extensionDependencies"`
	GalleryBanner         struct {
		Color string `json:"color"`
		Theme string `json:"theme"`
	} `json:"galleryBanner"`
	Author struct {
		Name string `json:"name"`
	} `json:"author"`
	Engines struct {
		Node   string `json:"node"`
		VSCode string `json:"vscode"`
	} `json:"engines"`
	Extension struct {
		ExtensionType string `json:"type"`
	} `json:"extension"`
	Bundles struct {
		Server  string `json:"server"`
		Browser string `json:"browser"`
	} `json:"bundles"`
	Scripts struct {
		Build   string `json:"cdebase:build"`
		Publish string `json:"cdebase:publish"`
	} `json:"scripts"`
}

type ExtensionAsset struct {
	Type    string
	Content string
}

type Manifest struct {
	NpmPackage

	Bundle string           `json:"bundle"`
	Assets []ExtensionAsset `json:"assets"`
}

func (m *Manifest) String() string {
	var str string

	data, _ := json.Marshal(m)
	str = string(data)

	return str
}

func (m *Manifest) ReadBundle(dir string) error {
	if m.Extension.ExtensionType != TYPE_SIMPLE {
		return nil
	}

	bundle, err := readFile(dir, m.Main)

	if err != nil {
		log.Fatal(err)
		return err
	}

	m.Bundle = string(bundle)
	return nil
}

func (m *Manifest) ReadAssets(dir string) error {
	assets := []ExtensionAsset{}
	m.ExtensionID = fmt.Sprintf("%s/%s", m.Publisher, m.Name)

	for asset, files := range EXTENSION_ASSETS {
		content := findFile(dir, files)
		if content != "" {
			assets = append(assets, ExtensionAsset{
				Type:    asset,
				Content: content,
			})
		}
	}

	m.Assets = assets

	return nil
}

func (m *Manifest) Validate() (bool, error) {
	if m.Name == "" && m.Publisher == "" {
		return false, errors.New(`extension manifest must contain "name" and "publisher" string properties (the extension ID is of the form "publisher/name" and uses these values)`)
	}

	if m.Name == "" {
		return false, fmt.Errorf(`extension manifest must contain a "name" string property for the extension name (the extension ID will be %q)`, m.Publisher+"/name")
	}

	if m.Publisher == "" {
		return false, fmt.Errorf(`extension manifest must contain a "publisher" string property referring to a username or organization name on Sourcegraph (the extension ID will be %q)`, "publisher/"+m.Name)
	}

	if m.ExtensionID == "" {
		return false, fmt.Errorf(`extension manifest must contain a "extensionId" (the extension ID will be %q)`, "publisher/"+m.ExtensionID)
	}

	return true, nil
}
