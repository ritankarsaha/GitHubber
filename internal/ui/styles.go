/*
 * GitHubber - UI Styling System
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Terminal styling and theme system using lipgloss
 */

package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	primaryColor   = lipgloss.Color("#00D9FF")
	secondaryColor = lipgloss.Color("#FF6B35")
	accentColor    = lipgloss.Color("#10D869")
	warningColor   = lipgloss.Color("#FFD23F")
	errorColor     = lipgloss.Color("#FF3366")
	textColor      = lipgloss.Color("#FFFFFF")
	mutedColor     = lipgloss.Color("#8B949E")
	bgColor        = lipgloss.Color("#0D1117")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(bgColor)

	// Header styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginBottom(1)

	// Menu styles
	MenuHeaderStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(0).
			PaddingLeft(1)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2).
			MarginBottom(0)

	MenuItemNumberStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			PaddingLeft(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			PaddingLeft(1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			PaddingLeft(1)

	InfoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			PaddingLeft(1)

	// Input styles
	PromptStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(0)

	InputStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			PaddingLeft(1).
			PaddingRight(1)

	// Content styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	CodeStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(lipgloss.Color("#161B22")).
			Padding(0, 1)

	// Repository info styles
	RepoInfoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			MarginBottom(2)

	RepoLabelStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	RepoValueStyle = lipgloss.NewStyle().
			Foreground(textColor)
)

// Themed emoji and icons
const (
	IconRepository = "📂"
	IconBranch     = "🌿"
	IconCommit     = "💾"
	IconRemote     = "🔄"
	IconHistory    = "📜"
	IconStash      = "📦"
	IconTag        = "🏷️"
	IconSuccess    = "✅"
	IconError      = "❌"
	IconWarning    = "⚠️"
	IconInfo       = "ℹ️"
	IconTool       = "🛠"
	IconGitHub     = "🐙"
	IconConfig     = "⚙️"
	IconExit       = "👋"
)

// Helper functions for consistent formatting
func FormatTitle(text string) string {
	return TitleStyle.Render("🚀 " + text)
}

func FormatSubtitle(text string) string {
	return SubtitleStyle.Render(text)
}

func FormatMenuHeader(icon, text string) string {
	return MenuHeaderStyle.Render("\n" + icon + " " + text + ":")
}

func FormatMenuItem(number int, text string) string {
	numStr := MenuItemNumberStyle.Render(fmt.Sprintf("%2d", number))
	return MenuItemStyle.Render(numStr + ". " + text)
}

func FormatSuccess(text string) string {
	return SuccessStyle.Render(IconSuccess + " " + text)
}

func FormatError(text string) string {
	return ErrorStyle.Render(IconError + " " + text)
}

func FormatWarning(text string) string {
	return WarningStyle.Render(IconWarning + " " + text)
}

func FormatInfo(text string) string {
	return InfoStyle.Render(IconInfo + " " + text)
}

func FormatPrompt(text string) string {
	return PromptStyle.Render("🎯 " + text)
}

func FormatRepoInfo(url, branch string) string {
	return RepoInfoStyle.Render(
		RepoLabelStyle.Render(IconRepository+" Repository: ") +
			RepoValueStyle.Render(url) + "\n" +
			RepoLabelStyle.Render(IconBranch+" Current Branch: ") +
			RepoValueStyle.Render(branch),
	)
}

func FormatBox(content string) string {
	return BoxStyle.Render(content)
}

func FormatCode(content string) string {
	return CodeStyle.Render(content)
}