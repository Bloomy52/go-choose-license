package ui

import "github.com/charmbracelet/lipgloss"

// Colors palette - Vibrant, modern, fun yet professional
var (
	ColorPrimary   = lipgloss.Color("#7D56F4") // Purple
	ColorSecondary = lipgloss.Color("#00D2FF") // Cyan
	ColorAccent    = lipgloss.Color("#FF79C6") // Pink / Magenta
	ColorSuccess   = lipgloss.Color("#50FA7B") // Green
	ColorWarning   = lipgloss.Color("#FFB86C") // Orange
	ColorSubtle    = lipgloss.Color("#6272A4") // Muted Blue/Gray
	ColorBgDark    = lipgloss.Color("#21222C") // Dark background
	ColorFgLight   = lipgloss.Color("#F8F8F2") // Crisp White
	ColorCardBg    = lipgloss.Color("#282A36") // Card Background
)

// Styles
var (
	// Main Container
	AppContainerStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Margin(0, 0)

	// Header Banner
	BannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	BannerBadge = lipgloss.NewStyle().
			Bold(true).
			Background(ColorPrimary).
			Foreground(ColorFgLight).
			Padding(0, 1).
			MarginRight(1)

	BannerSubtitle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Italic(true)

	// Progress & Breadcrumb
	ProgressStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
			MarginBottom(1)

	CategoryBadge = lipgloss.NewStyle().
			Bold(true).
			Background(ColorSubtle).
			Foreground(ColorFgLight).
			Padding(0, 1).
			MarginRight(1)

	// Cards & Containers
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginBottom(1)

	QuestionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorFgLight).
				MarginBottom(1)

	QuestionSubtitleStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				MarginBottom(1)

	// Options
	OptionItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorFgLight)

	OptionSelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Bold(true).
				Foreground(ColorSecondary).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(ColorSecondary)

	OptionDescStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Italic(true).
			PaddingLeft(4)

	// Tags & Badges
	PermissiveBadge = lipgloss.NewStyle().
			Bold(true).
			Background(ColorSuccess).
			Foreground(ColorBgDark).
			Padding(0, 1)

	CopyleftBadge = lipgloss.NewStyle().
			Bold(true).
			Background(ColorWarning).
			Foreground(ColorBgDark).
			Padding(0, 1)

	// Help Bar
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Inputs
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1)

	// Output/Success Box
	SuccessCardStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(ColorSuccess).
				Padding(1, 2).
				MarginTop(1)
)
