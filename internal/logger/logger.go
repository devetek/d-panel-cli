package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Success(msg string) {
	fmt.Println(lipgloss.NewStyle().Bold(true).Render("✅ " + msg))
}

func Error(msg string) {
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ca1414ff")).Render("❌ " + msg))
}
