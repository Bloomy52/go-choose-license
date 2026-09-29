package ui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Bloomy52/go-choose-license/internal/license"
	"github.com/Bloomy52/go-choose-license/internal/ui"
)

func TestQuestionnaireNavigation(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	model := ui.InitialModel(reg)

	// Verify main menu renders header
	viewStr := model.View()
	if !strings.Contains(viewStr, "GO CHOOSE YOUR LICENSE") {
		t.Errorf("Expected view to contain banner header")
	}

	// Press Enter on menu item 0 (Interactive Questionnaire)
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := m.(ui.Model)

	viewQ1 := m2.View()
	if !strings.Contains(viewQ1, "Question 1 of 5") {
		t.Errorf("Expected Question 1 of 5 in view, got: %s", viewQ1)
	}
	if !strings.Contains(viewQ1, "closed-source") {
		t.Errorf("Expected Q1 title in view")
	}

	// Press Enter on Q1 option 0 (Yes, allow proprietary) -> leads to Q2
	m, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := m.(ui.Model)

	viewQ2 := m3.View()
	if !strings.Contains(viewQ2, "public domain") {
		t.Errorf("Expected Q2 title in view, got: %s", viewQ2)
	}

	// Press 'b' to go back to Q1
	m, _ = m3.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m4 := m.(ui.Model)

	viewQ1Back := m4.View()
	if !strings.Contains(viewQ1Back, "Question 1 of 5") {
		t.Errorf("Expected to return to Question 1 after pressing 'b'")
	}
}

func TestLanguageSelectFilter(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	model := ui.InitialModel(reg)

	// Move cursor to option 1 (Language Norms)
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, _ = m.(ui.Model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	mLang := m.(ui.Model)

	viewLang := mLang.View()
	if !strings.Contains(viewLang, "Programming Language Norms") {
		t.Errorf("Expected Language Select title in view")
	}
	if !strings.Contains(viewLang, "Python") || !strings.Contains(viewLang, "Go") {
		t.Errorf("Expected languages listed in view")
	}
}

func TestResultCardMaxWidth(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	model := ui.InitialModel(reg)

	// Update window size with width 60
	m, _ := model.Update(tea.WindowSizeMsg{Width: 60, Height: 24})
	mRes := m.(ui.Model)

	viewStr := mRes.View()
	lines := strings.Split(viewStr, "\n")
	for _, l := range lines {
		w := lipgloss.Width(l)
		if w > 60 {
			t.Errorf("Line width %d exceeds target width 60 in line: %q", w, l)
		}
	}
}

func TestPackageManagerSelection(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatal(err)
	}
	for _, norm := range license.GetPackageManagerNorms() {
		t.Run(norm.Language, func(t *testing.T) {
			var model tea.Model = ui.InitialModel(reg)
			press := func(key tea.KeyMsg) { model, _ = model.Update(key) }
			press(tea.KeyMsg{Type: tea.KeyDown})
			press(tea.KeyMsg{Type: tea.KeyDown})
			press(tea.KeyMsg{Type: tea.KeyEnter})
			press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(norm.Language)})
			if !strings.Contains(model.View(), norm.Language) {
				t.Fatal("Package manager missing from filtered list")
			}
			press(tea.KeyMsg{Type: tea.KeyEnter})
			view := model.View()
			if !strings.Contains(view, "License Recommendation") || !strings.Contains(view, "Community norm for "+norm.Language) {
				t.Fatalf("Missing recommendation or community note: %s", view)
			}
			for _, id := range norm.LicenseIDs {
				lic, ok := reg.Get(id)
				if !ok || !strings.Contains(view, lic.Name) {
					t.Errorf("Missing recommended license %s", id)
				}
			}
			press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")})
			press(tea.KeyMsg{Type: tea.KeyEsc})
			press(tea.KeyMsg{Type: tea.KeyEnter})
			if !strings.Contains(model.View(), "Generate") {
				t.Fatal("Cannot reach license generation")
			}
			press(tea.KeyMsg{Type: tea.KeyEsc})
			press(tea.KeyMsg{Type: tea.KeyEsc})
			if !strings.Contains(model.View(), "Package Manager Norms") || !strings.Contains(model.View(), norm.Language) {
				t.Fatal("Back navigation lost package manager selection")
			}
			press(tea.KeyMsg{Type: tea.KeyEsc})
			press(tea.KeyMsg{Type: tea.KeyUp})
			press(tea.KeyMsg{Type: tea.KeyEnter})
			if !strings.Contains(model.View(), "Python") {
				t.Fatal("Package manager search leaked into language selection")
			}
		})
	}
}

func TestPackageManagerEmptySearch(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var model tea.Model = ui.InitialModel(reg)
	for _, key := range []tea.KeyMsg{
		{Type: tea.KeyDown}, {Type: tea.KeyDown}, {Type: tea.KeyEnter},
		{Type: tea.KeyDown}, {Type: tea.KeyRunes, Runes: []rune("no-such-package-manager")},
		{Type: tea.KeyEnter},
	} {
		model, _ = model.Update(key)
	}
	if !strings.Contains(model.View(), "No matching package managers found.") {
		t.Fatal("Expected empty search state")
	}
}
