package ui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go-choose-license/internal/license"
	"go-choose-license/internal/ui"
)

func TestQuestionnaireNavigation(t *testing.T) {
	reg, err := license.LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	model := ui.InitialModel(reg)

	// Verify main menu renders header
	viewStr := model.View()
	if !strings.Contains(viewStr, "GO LICENSE CHOOSER") {
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
