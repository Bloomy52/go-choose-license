package license

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Embedded FS containing license definition YAML and raw license text files.
//
//go:embed embedded/licenses.yaml embedded/licenses/*.txt
var embeddedFS embed.FS

type License struct {
	ID         string `yaml:"id"`
	Name       string `yaml:"name"`
	File       string `yaml:"file"`
	Permissive bool   `yaml:"permissive"`
	Summary    string `yaml:"-"`
	Text       string `yaml:"-"`
}

type LicenseFile struct {
	Licenses []License `yaml:"licenses"`
}

type Registry struct {
	licenses map[string]License
	list     []License
}

var GlobalRegistry *Registry

func LoadRegistry() (*Registry, error) {
	data, err := embeddedFS.ReadFile("embedded/licenses.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded licenses.yaml: %w", err)
	}

	var lf LicenseFile
	if err := yaml.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("failed to parse licenses.yaml: %w", err)
	}

	reg := &Registry{
		licenses: make(map[string]License),
		list:     lf.Licenses,
	}

	// Summaries from planning docs
	summaries := map[string]string{
		"mit":          "Do anything you want, just keep the copyright notice.",
		"apache-2.0":   "Permissive, business-friendly license with an explicit patent grant.",
		"bsd-3-clause": "Permissive like MIT, but you can't use the authors' names to endorse your product.",
		"bsd-2-clause": "Permissive license, just keep the copyright notice.",
		"unlicense":    "Public domain — no rights reserved, no attribution needed.",
		"bsd0":         "Public domain-style license with zero conditions.",
		"lgpl-2.1":     "You can link to this library freely, but changes to the library itself must be shared.",
		"lgpl-3.0":     "Like LGPL-2.1, but updated to match GPLv3's terms.",
		"mpl-2.0":      "Changes to MPL-licensed files must be shared, but you can mix it with proprietary code.",
		"epl-2.0":      "File-level copyleft with a patent grant, common in the Eclipse ecosystem.",
		"gpl-2.0":      "Any code you distribute based on this must also be open-sourced under GPL-2.0.",
		"gpl-3.0":      "Like GPL-2.0, but adds patent protection and blocks hardware lockdown (Tivoization).",
		"agpl-3.0":     "Like GPL-3.0, but also covers software used over a network, like SaaS.",
		"artistic-2.0": "Permissive license popular in Perl, with some copyleft-like conditions on modifications.",
		"postgresql":   "Permissive license nearly identical to MIT or BSD.",
	}

	for i, lic := range lf.Licenses {
		pathInFS := filepath.Join("embedded", lic.File)
		textBytes, err := embeddedFS.ReadFile(pathInFS)
		if err == nil {
			lic.Text = string(textBytes)
		} else {
			lic.Text = fmt.Sprintf("License text for %s could not be loaded: %v", lic.Name, err)
		}

		if sum, ok := summaries[lic.ID]; ok {
			lic.Summary = sum
		} else {
			if lic.Permissive {
				lic.Summary = "Permissive Open Source License"
			} else {
				lic.Summary = "Copyleft Open Source License"
			}
		}

		reg.licenses[lic.ID] = lic
		reg.list[i] = lic
	}

	GlobalRegistry = reg
	return reg, nil
}

func (r *Registry) Get(id string) (License, bool) {
	lic, ok := r.licenses[strings.ToLower(id)]
	return lic, ok
}

func (r *Registry) List() []License {
	return r.list
}
