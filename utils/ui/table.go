package ui

import (
	"strings"
	"unicode/utf8"
)

func DrawTable(title string, content []string) string {
	if len(content) == 0 {
		return ""
	}

	maxW := utf8.RuneCountInString(title) + 4
	for _, line := range content {
		cLen := utf8.RuneCountInString(line)
		if cLen > maxW {
			maxW = cLen
		}
	}
	maxW += 2

	var sb strings.Builder

	titleLen := utf8.RuneCountInString(title)
	sideBarLen := (maxW - titleLen - 2) / 2
	if sideBarLen < 1 {
		sideBarLen = 1
	}

	sb.WriteString("╭" + strings.Repeat("─", sideBarLen))
	sb.WriteString(" " + title + " ")
	rightBarLen := maxW - sideBarLen - titleLen - 2
	sb.WriteString(strings.Repeat("─", rightBarLen) + "╮\n")

	for _, line := range content {
		lineLen := utf8.RuneCountInString(line)
		padding := maxW - lineLen
		sb.WriteString("│ " + line + strings.Repeat(" ", padding-1) + "│\n")
	}

	sb.WriteString("╰" + strings.Repeat("─", maxW) + "╯")

	return sb.String()
}
