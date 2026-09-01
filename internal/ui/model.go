package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go-choose-license/internal/license"
)

type State int

const (
	StateMenu State = iota
	StateQuestionnaire
	StateLanguageSelect
	StateCatalog
	StateResult
	StateViewText
	StateGenerate
	StateSuccess
)

type Model struct {
	state             State
	registry          *license.Registry
	questions         map[license.QuestionID]license.Question
	currentQuestionID license.QuestionID
	historyStack      []license.QuestionID

	// Selection cursor
	menuCursor   int
	optCursor    int
	langCursor   int
	catCursor    int
	resultCursor int

	// Data selections
	selectedLanguage license.LanguageNorm
	recommendations  []license.License
	chosenLicense    license.License
	prevState        State

	// Inputs & Viewport
	searchInput textinput.Model
	yearInput   textinput.Model
	authorInput textinput.Model
	activeInput int // 0: year, 1: author
	viewport    viewport.Model

	// Results / Stats
	generatedFilePath string
	errMessage        string

	// Window dimensions
	width  int
	height int
}

func InitialModel(reg *license.Registry) Model {
	si := textinput.New()
	si.Placeholder = "Search programming language..."
	si.CharLimit = 50

	yi := textinput.New()
	yi.Placeholder = fmt.Sprintf("%d", time.Now().Year())
	yi.SetValue(fmt.Sprintf("%d", time.Now().Year()))
	yi.CharLimit = 10

	ai := textinput.New()
	ai.Placeholder = "Your Name or Organization"
	ai.SetValue(getGitUserName())
	ai.CharLimit = 100

	vp := viewport.New(80, 20)

	return Model{
		state:             StateMenu,
		registry:          reg,
		questions:         license.GetQuestionsMap(),
		currentQuestionID: license.Q1,
		historyStack:      make([]license.QuestionID, 0),
		searchInput:       si,
		yearInput:         yi,
		authorInput:       ai,
		viewport:          vp,
		width:             80,
		height:            24,
	}
}

func getGitUserName() string {
	out, err := exec.Command("git", "config", "user.name").Output()
	if err == nil {
		name := strings.TrimSpace(string(out))
		if name != "" {
			return name
		}
	}
	return "Your Name"
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		vpWidth := msg.Width - 10
		if vpWidth < 10 {
			vpWidth = 10
		}
		vpHeight := msg.Height - 10
		if vpHeight < 5 {
			vpHeight = 5
		}
		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateMenu:
		m, cmd = m.updateMenu(msg)
	case StateQuestionnaire:
		m, cmd = m.updateQuestionnaire(msg)
	case StateLanguageSelect:
		m, cmd = m.updateLanguageSelect(msg)
	case StateCatalog:
		m, cmd = m.updateCatalog(msg)
	case StateResult:
		m, cmd = m.updateResult(msg)
	case StateViewText:
		m, cmd = m.updateViewText(msg)
	case StateGenerate:
		m, cmd = m.updateGenerate(msg)
	case StateSuccess:
		m, cmd = m.updateSuccess(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Menu Update
func (m Model) updateMenu(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "j":
			if m.menuCursor < 3 {
				m.menuCursor++
			}
		case "enter", "space":
			switch m.menuCursor {
			case 0: // Questionnaire
				m.state = StateQuestionnaire
				m.currentQuestionID = license.Q1
				m.historyStack = nil
				m.optCursor = 0
			case 1: // Language
				m.state = StateLanguageSelect
				m.langCursor = 0
				m.searchInput.Focus()
				return m, textinput.Blink
			case 2: // View All Licenses
				m.state = StateCatalog
				m.catCursor = 0
			case 3: // Quit
				return m, tea.Quit
			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// Questionnaire Update
func (m Model) updateQuestionnaire(msg tea.Msg) (Model, tea.Cmd) {
	q := m.questions[m.currentQuestionID]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.optCursor > 0 {
				m.optCursor--
			}
		case "down", "j":
			if m.optCursor < len(q.Options)-1 {
				m.optCursor++
			}
		case "b", "backspace":
			if len(m.historyStack) > 0 {
				m.currentQuestionID = m.historyStack[len(m.historyStack)-1]
				m.historyStack = m.historyStack[:len(m.historyStack)-1]
				m.optCursor = 0
			} else {
				m.state = StateMenu
			}
		case "enter", "space":
			opt := q.Options[m.optCursor]
			if len(opt.LicenseIDs) > 0 {
				// We reached a license recommendation!
				m.recommendations = nil
				for _, id := range opt.LicenseIDs {
					if lic, ok := m.registry.Get(id); ok {
						m.recommendations = append(m.recommendations, lic)
					}
				}
				if len(m.recommendations) > 0 {
					m.chosenLicense = m.recommendations[0]
					m.resultCursor = 0
					m.prevState = StateQuestionnaire
					m.state = StateResult
				}
			} else if opt.NextQuestion != "" {
				m.historyStack = append(m.historyStack, m.currentQuestionID)
				m.currentQuestionID = opt.NextQuestion
				m.optCursor = 0
			}
		case "q":
			m.state = StateMenu
		}
	}
	return m, nil
}

// Language Select Update
func (m Model) updateLanguageSelect(msg tea.Msg) (Model, tea.Cmd) {
	filteredLangs := m.getFilteredLanguages()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if m.langCursor > 0 {
				m.langCursor--
			}
			return m, nil
		case "down", "ctrl+n":
			if m.langCursor < len(filteredLangs)-1 {
				m.langCursor++
			}
			return m, nil
		case "enter":
			if len(filteredLangs) > 0 && m.langCursor < len(filteredLangs) {
				m.selectedLanguage = filteredLangs[m.langCursor]
				m.recommendations = nil
				for _, id := range m.selectedLanguage.LicenseIDs {
					if lic, ok := m.registry.Get(id); ok {
						m.recommendations = append(m.recommendations, lic)
					}
				}
				if len(m.recommendations) > 0 {
					m.chosenLicense = m.recommendations[0]
					m.resultCursor = 0
					m.prevState = StateLanguageSelect
					m.state = StateResult
				}
			}
			return m, nil
		case "esc", "ctrl+b":
			m.state = StateMenu
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	// reset cursor if input changed
	if m.langCursor >= len(m.getFilteredLanguages()) {
		m.langCursor = 0
	}
	return m, cmd
}

func (m Model) getFilteredLanguages() []license.LanguageNorm {
	all := license.GetLanguageNorms()
	query := strings.TrimSpace(strings.ToLower(m.searchInput.Value()))
	if query == "" {
		return all
	}
	var filtered []license.LanguageNorm
	for _, l := range all {
		if strings.Contains(strings.ToLower(l.Language), query) || strings.Contains(strings.ToLower(l.Note), query) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// Catalog Update
func (m Model) updateCatalog(msg tea.Msg) (Model, tea.Cmd) {
	allLics := m.registry.List()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.catCursor > 0 {
				m.catCursor--
			}
		case "down", "j":
			if m.catCursor < len(allLics)-1 {
				m.catCursor++
			}
		case "enter", "space":
			if m.catCursor < len(allLics) {
				m.chosenLicense = allLics[m.catCursor]
				m.recommendations = []license.License{m.chosenLicense}
				m.resultCursor = 0
				m.prevState = StateCatalog
				m.state = StateResult
			}
		case "b", "esc", "q":
			m.state = StateMenu
		}
	}
	return m, nil
}

// Result Update
func (m Model) updateResult(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if len(m.recommendations) > 1 && m.resultCursor > 0 {
				m.resultCursor--
				m.chosenLicense = m.recommendations[m.resultCursor]
			}
		case "down", "j":
			if len(m.recommendations) > 1 && m.resultCursor < len(m.recommendations)-1 {
				m.resultCursor++
				m.chosenLicense = m.recommendations[m.resultCursor]
			}
		case "v", "V":
			vpWidth := m.viewport.Width
			if vpWidth <= 0 {
				vpWidth = 70
			}
			wrapped := lipgloss.NewStyle().MaxWidth(vpWidth).Render(m.chosenLicense.Text)
			m.viewport.SetContent(wrapped)
			m.viewport.GotoTop()
			m.state = StateViewText
		case "g", "G", "enter":
			m.state = StateGenerate
			m.activeInput = 0
			m.yearInput.Focus()
			m.authorInput.Blur()
			return m, textinput.Blink
		case "b", "esc":
			m.state = m.prevState
		case "q":
			m.state = StateMenu
		}
	}
	return m, nil
}

// View Text Update
func (m Model) updateViewText(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b", "esc", "q":
			m.state = StateResult
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// Generate Update
func (m Model) updateGenerate(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "up":
			if m.activeInput == 0 {
				m.activeInput = 1
				m.yearInput.Blur()
				m.authorInput.Focus()
			} else {
				m.activeInput = 0
				m.authorInput.Blur()
				m.yearInput.Focus()
			}
			return m, textinput.Blink
		case "esc", "b":
			m.state = StateResult
			return m, nil
		case "enter":
			err := m.saveLicenseFile()
			if err != nil {
				m.errMessage = err.Error()
			} else {
				m.errMessage = ""
				m.state = StateSuccess
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.activeInput == 0 {
		m.yearInput, cmd = m.yearInput.Update(msg)
	} else {
		m.authorInput, cmd = m.authorInput.Update(msg)
	}
	return m, cmd
}

func (m *Model) saveLicenseFile() error {
	yearStr := strings.TrimSpace(m.yearInput.Value())
	if yearStr == "" {
		yearStr = fmt.Sprintf("%d", time.Now().Year())
	}
	authorStr := strings.TrimSpace(m.authorInput.Value())
	if authorStr == "" {
		authorStr = "Copyright Owner"
	}

	text := m.chosenLicense.Text

	// Standard template replacements
	text = strings.ReplaceAll(text, "<year>", yearStr)
	text = strings.ReplaceAll(text, "[year]", yearStr)
	text = strings.ReplaceAll(text, "[yyyy]", yearStr)
	text = strings.ReplaceAll(text, "<copyright holders>", authorStr)
	text = strings.ReplaceAll(text, "[fullname]", authorStr)
	text = strings.ReplaceAll(text, "[name of copyright owner]", authorStr)
	text = strings.ReplaceAll(text, "<name of author>", authorStr)

	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	targetPath := filepath.Join(pwd, "LICENSE")

	err = os.WriteFile(targetPath, []byte(text), 0644)
	if err != nil {
		return fmt.Errorf("failed to write LICENSE file: %w", err)
	}

	m.generatedFilePath = targetPath
	return nil
}

// Success Update
func (m Model) updateSuccess(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "space", "q", "esc":
			m.state = StateMenu
		}
	}
	return m, nil
}

func (m Model) renderAppContainer(body string) string {
	if m.width > 0 {
		return AppContainerStyle.MaxWidth(m.width).Render(body)
	}
	return AppContainerStyle.Render(body)
}

func (m Model) renderCard(content string) string {
	if m.width > 0 {
		maxW := m.width - 4
		if maxW < 20 {
			maxW = 20
		}
		return CardStyle.MaxWidth(maxW).Render(content)
	}
	return CardStyle.Render(content)
}

func (m Model) renderSuccessCard(content string) string {
	if m.width > 0 {
		maxW := m.width - 4
		if maxW < 20 {
			maxW = 20
		}
		return SuccessCardStyle.MaxWidth(maxW).Render(content)
	}
	return SuccessCardStyle.Render(content)
}

// View Functions
func (m Model) View() string {
	var body string

	switch m.state {
	case StateMenu:
		body = m.viewMenu()
	case StateQuestionnaire:
		body = m.viewQuestionnaire()
	case StateLanguageSelect:
		body = m.viewLanguageSelect()
	case StateCatalog:
		body = m.viewCatalog()
	case StateResult:
		body = m.viewResult()
	case StateViewText:
		body = m.viewViewText()
	case StateGenerate:
		body = m.viewGenerate()
	case StateSuccess:
		body = m.viewSuccess()
	}

	return m.renderAppContainer(body)
}

func (m Model) renderHeader(title string, subtitle string) string {
	banner := BannerBadge.Render("GO CHOOSE YOUR LICENSE") + BannerStyle.Render(title)
	if subtitle != "" {
		banner += "\n" + BannerSubtitle.Render(subtitle)
	}
	return banner + "\n\n"
}

func (m Model) viewMenu() string {
	h := m.renderHeader("Main Menu", "Select how you would like to choose an open source license for your project")

	options := []struct {
		title string
		desc  string
	}{
		{"1. Interactive Questionnaire", "Answer a few questions to find the license matching your goals"},
		{"2. Programming Language Norms", "Choose a license based on community conventions for your language"},
		{"3. Browse All Licenses", "View full catalog of available permissive & copyleft licenses"},
		{"4. Exit", "Quit application"},
	}

	var sb strings.Builder
	for i, opt := range options {
		if i == m.menuCursor {
			sb.WriteString(OptionSelectedItemStyle.Render(opt.title) + "\n")
			sb.WriteString(OptionDescStyle.Render(opt.desc) + "\n\n")
		} else {
			sb.WriteString(OptionItemStyle.Render(opt.title) + "\n")
			sb.WriteString(OptionDescStyle.Render(opt.desc) + "\n\n")
		}
	}

	card := m.renderCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("%s move • %s select • %s quit", KeyStyle.Render("↑/↓ / j/k"), KeyStyle.Render("Enter"), KeyStyle.Render("q")))

	return h + card + help
}

func (m Model) viewQuestionnaire() string {
	q := m.questions[m.currentQuestionID]

	categoryBadge := CategoryBadge.Render(strings.ToUpper(q.Category))
	stepStr := ProgressStyle.Render(fmt.Sprintf("Question %d of %d", q.StepNum, q.Total))
	h := m.renderHeader("Guided Questionnaire", "")

	var sb strings.Builder
	sb.WriteString(categoryBadge + " " + stepStr + "\n\n")
	sb.WriteString(QuestionTitleStyle.Render(q.Title) + "\n")
	if q.Subtitle != "" {
		sb.WriteString(QuestionSubtitleStyle.Render(q.Subtitle) + "\n\n")
	}

	for i, opt := range q.Options {
		if i == m.optCursor {
			sb.WriteString(OptionSelectedItemStyle.Render(opt.Text) + "\n")
			if opt.Desc != "" {
				sb.WriteString(OptionDescStyle.Render(opt.Desc) + "\n\n")
			}
		} else {
			sb.WriteString(OptionItemStyle.Render(opt.Text) + "\n")
			if opt.Desc != "" {
				sb.WriteString(OptionDescStyle.Render(opt.Desc) + "\n\n")
			}
		}
	}

	card := m.renderCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("%s navigate • %s confirm • %s back • %s main menu", KeyStyle.Render("↑/↓"), KeyStyle.Render("Enter"), KeyStyle.Render("b / Backspace"), KeyStyle.Render("q")))

	return h + card + help
}

func (m Model) viewLanguageSelect() string {
	h := m.renderHeader("Programming Language Norms", "Filter and select your primary programming language")

	inputView := InputStyle.Render(m.searchInput.View())

	filtered := m.getFilteredLanguages()

	var sb strings.Builder
	sb.WriteString(inputView + "\n\n")

	if len(filtered) == 0 {
		sb.WriteString(OptionDescStyle.Render("No matching languages found.") + "\n")
	} else {
		maxVisible := 8
		start := 0
		if m.langCursor >= maxVisible {
			start = m.langCursor - maxVisible + 1
		}
		end := start + maxVisible
		if end > len(filtered) {
			end = len(filtered)
		}

		for i := start; i < end; i++ {
			lang := filtered[i]
			licenseNames := strings.Join(lang.LicenseIDs, " / ")
			line := fmt.Sprintf("%-18s → Recommended: %s", lang.Language, strings.ToUpper(licenseNames))
			if i == m.langCursor {
				sb.WriteString(OptionSelectedItemStyle.Render(line) + "\n")
				sb.WriteString(OptionDescStyle.Render(lang.Note) + "\n")
			} else {
				sb.WriteString(OptionItemStyle.Render(line) + "\n")
			}
		}
	}

	card := m.renderCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("%s filter • %s/%s move • %s choose • %s menu", KeyStyle.Render("Type"), KeyStyle.Render("↑/↓"), KeyStyle.Render("Ctrl+n/p"), KeyStyle.Render("Enter"), KeyStyle.Render("Esc")))

	return h + card + help
}

func (m Model) viewCatalog() string {
	h := m.renderHeader("License Catalog", "Select any license to view details and generate a LICENSE file")

	all := m.registry.List()
	var sb strings.Builder

	for i, lic := range all {
		var badge string
		if lic.Permissive {
			badge = PermissiveBadge.Render("PERMISSIVE")
		} else {
			badge = CopyleftBadge.Render("COPYLEFT")
		}

		line := fmt.Sprintf("%-40s %s", lic.Name, badge)
		if i == m.catCursor {
			sb.WriteString(OptionSelectedItemStyle.Render(line) + "\n")
			sb.WriteString(OptionDescStyle.Render(lic.Summary) + "\n\n")
		} else {
			sb.WriteString(OptionItemStyle.Render(line) + "\n")
		}
	}

	card := m.renderCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("%s move • %s select • %s main menu", KeyStyle.Render("↑/↓ / j/k"), KeyStyle.Render("Enter"), KeyStyle.Render("b / Esc")))

	return h + card + help
}

func (m Model) viewResult() string {
	h := m.renderHeader("License Recommendation", "Based on your criteria, here is the suggested license")

	var sb strings.Builder

	if len(m.recommendations) > 1 {
		sb.WriteString(QuestionSubtitleStyle.Render("Multiple matches available! Select your preferred option:") + "\n\n")
		for i, lic := range m.recommendations {
			var badge string
			if lic.Permissive {
				badge = PermissiveBadge.Render("PERMISSIVE")
			} else {
				badge = CopyleftBadge.Render("COPYLEFT")
			}
			line := fmt.Sprintf("%s (%s)", lic.Name, lic.ID)
			if i == m.resultCursor {
				sb.WriteString(OptionSelectedItemStyle.Render(line+" ") + badge + "\n")
				sb.WriteString(OptionDescStyle.Render(lic.Summary) + "\n\n")
			} else {
				sb.WriteString(OptionItemStyle.Render(line+" ") + badge + "\n\n")
			}
		}
	} else {
		lic := m.chosenLicense
		var badge string
		if lic.Permissive {
			badge = PermissiveBadge.Render("PERMISSIVE")
		} else {
			badge = CopyleftBadge.Render("COPYLEFT")
		}

		sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render(lic.Name) + " " + badge + "\n\n")
		sb.WriteString(QuestionTitleStyle.Render("Summary:") + " " + lic.Summary + "\n\n")
		if m.selectedLanguage.Language != "" {
			sb.WriteString(OptionDescStyle.Render(fmt.Sprintf("Community norm for %s: %s", m.selectedLanguage.Language, m.selectedLanguage.Note)) + "\n\n")
		}
	}

	card := m.renderCard(sb.String())

	actions := lipgloss.NewStyle().Foreground(ColorFgLight).Render(fmt.Sprintf(
		"\nActions:\n  %s Generate LICENSE file\n  %s View full license text\n  %s Go back\n  %s Main menu",
		KeyStyle.Render("[G / Enter]"),
		KeyStyle.Render("[V]"),
		KeyStyle.Render("[B]"),
		KeyStyle.Render("[Q]"),
	))

	return h + card + actions
}

func (m Model) viewViewText() string {
	h := m.renderHeader(fmt.Sprintf("Full License Text — %s", m.chosenLicense.Name), "Scroll with arrow keys or PgUp/PgDn")

	content := m.renderCard(m.viewport.View())
	help := HelpStyle.Render(fmt.Sprintf("%s scroll • %s back to recommendation", KeyStyle.Render("↑/↓ / PgUp/PgDn"), KeyStyle.Render("b / Esc / q")))

	return h + content + help
}

func (m Model) viewGenerate() string {
	h := m.renderHeader("Generate LICENSE File", fmt.Sprintf("Creating LICENSE file for %s", m.chosenLicense.Name))

	var sb strings.Builder
	sb.WriteString(QuestionTitleStyle.Render("Please provide copyright details:") + "\n\n")

	yearLabel := "Copyright Year:"
	if m.activeInput == 0 {
		yearLabel = OptionSelectedItemStyle.Render(yearLabel)
	} else {
		yearLabel = OptionItemStyle.Render(yearLabel)
	}
	sb.WriteString(yearLabel + "\n" + InputStyle.Render(m.yearInput.View()) + "\n\n")

	authorLabel := "Copyright Holder / Author Name:"
	if m.activeInput == 1 {
		authorLabel = OptionSelectedItemStyle.Render(authorLabel)
	} else {
		authorLabel = OptionItemStyle.Render(authorLabel)
	}
	sb.WriteString(authorLabel + "\n" + InputStyle.Render(m.authorInput.View()) + "\n\n")

	if m.errMessage != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render("Error: "+m.errMessage) + "\n")
	}

	card := m.renderCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("%s switch fields • %s write LICENSE file • %s cancel", KeyStyle.Render("Tab / ↑/↓"), KeyStyle.Render("Enter"), KeyStyle.Render("Esc")))

	return h + card + help
}

func (m Model) viewSuccess() string {
	h := m.renderHeader("License Generated Successfully! 🎉", "")

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render("✓ LICENSE file created:") + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorSecondary).Render(m.generatedFilePath) + "\n\n")
	sb.WriteString(QuestionSubtitleStyle.Render(fmt.Sprintf("License: %s", m.chosenLicense.Name)) + "\n")
	sb.WriteString(QuestionSubtitleStyle.Render(fmt.Sprintf("Year: %s | Author: %s", m.yearInput.Value(), m.authorInput.Value())) + "\n")

	card := m.renderSuccessCard(sb.String())
	help := HelpStyle.Render(fmt.Sprintf("Press %s or %s to return to main menu", KeyStyle.Render("Enter"), KeyStyle.Render("q")))

	return h + card + help
}
