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
		"mit":          "Simple & Popular License, Requires Copyright & Notice Preserved, Warranty Disclaimed",
		"apache-2.0":   "Popular for Businesses, Explicit Patent Grant, Copyright & Notice Preserved, Warranty Disclaimed, State Changes",
		"bsd-3-clause": "Popular for Academics, Non-Endorsement Clause, Copyright & Notice Preserved, Warranty Disclaimed",
		"bsd-2-clause": "Popular for Academics, Requires Copyright & Notice Preserved, Warranty Disclaimed",
		"unlicense":    "Public Domain Dedication, Warranty Disclaimed",
		"bsd0":         "No Conditions Present, Warranty Disclaimed",
		"lgpl-2.1":     "Used mainly for libraries, usage in libraries does not require release, modification requires release",
		"lgpl-3.0":     "Updated weak copyleft compatible with GPLv3",
		"mpl-2.0":      "File Level Copyleft",
		"epl-2.0":      "File Level Copyleft, popular in Eclipse ecosystem",
		"gpl-2.0":      "Source Code Modifications Released Under GPLv2, Tivoization Allowed",
		"gpl-3.0":      "Extension to GPLv2, Anti-Tivoization Clause, Explicit Patent Grant",
		"agpl-3.0":     "Removes Network Distribution/SaaS Source Code Distribution Loophole",
		"artistic-2.0": "Permissive/Copyleft hybrid popular in Perl ecosystem",
		"postgresql":   "Liberal permissive license similar to MIT/BSD",
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
