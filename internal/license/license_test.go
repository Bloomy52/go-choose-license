package license_test

import (
	"testing"

	"github.com/Bloomy52/go-choose-license/internal/license"
)

func TestLoadRegistry(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	list := reg.List()
	if len(list) < 15 {
		t.Errorf("Expected at least 15 licenses, got %d", len(list))
	}

	expectedIDs := []string{
		"mit", "apache-2.0", "bsd-3-clause", "bsd-2-clause", "unlicense",
		"bsd0", "lgpl-2.1", "lgpl-3.0", "mpl-2.0", "epl-2.0",
		"gpl-2.0", "gpl-3.0", "agpl-3.0", "artistic-2.0", "postgresql",
	}

	for _, id := range expectedIDs {
		lic, ok := reg.Get(id)
		if !ok {
			t.Errorf("Expected license ID %q to exist in registry", id)
			continue
		}
		if lic.Text == "" {
			t.Errorf("License text for %q should not be empty", id)
		}
	}
}

func TestGetQuestionsMap(t *testing.T) {
	questions := license.GetQuestionsMap()
	if len(questions) != 9 {
		t.Errorf("Expected 9 questions, got %d", len(questions))
	}

	q1, ok := questions[license.Q1]
	if !ok {
		t.Fatalf("Q1 missing from questions map")
	}
	if len(q1.Options) != 2 {
		t.Errorf("Q1 should have 2 options, got %d", len(q1.Options))
	}

	// Verify Q1 option 0 leads to Q2 (permissive) and option 1 leads to Q5 (copyleft)
	if q1.Options[0].NextQuestion != license.Q2 {
		t.Errorf("Expected Q1 Option 0 to lead to Q2, got %s", q1.Options[0].NextQuestion)
	}
	if q1.Options[1].NextQuestion != license.Q5 {
		t.Errorf("Expected Q1 Option 1 to lead to Q5, got %s", q1.Options[1].NextQuestion)
	}
}

func TestGetLanguageNorms(t *testing.T) {
	norms := license.GetLanguageNorms()
	if len(norms) < 35 {
		t.Errorf("Expected at least 35 language norms, got %d", len(norms))
	}

	// Check Go norm
	foundGo := false
	for _, n := range norms {
		if n.Language == "Go" {
			foundGo = true
			if len(n.LicenseIDs) == 0 || n.LicenseIDs[0] != "apache-2.0" {
				t.Errorf("Expected Go default license to be apache-2.0, got %v", n.LicenseIDs)
			}
		}
	}

	if !foundGo {
		t.Errorf("Language 'Go' not found in language norms")
	}
}
