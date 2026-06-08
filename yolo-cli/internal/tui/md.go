package tui

import (
	"strings"
)

// ANSI escape codes for terminal styling.
const (
	ansiBold      = "\033[1m"
	ansiItalic    = "\033[3m"
	ansiUnderline = "\033[4m"
	ansiReset     = "\033[0m"
	ansiFaint     = "\033[2m"
)

// Simple markdown to ANSI-styled terminal text renderer.
// Supports: **bold** *italic* `code` ```code blocks``` # headings - lists

func mdRender(content string, width int) string {
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	var out []string
	inCodeBlock := false
	var codeBuf []string

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// End code block
				out = append(out, renderCodeBlock(codeBuf, width))
				codeBuf = nil
				inCodeBlock = false
			} else {
				inCodeBlock = true
				codeBuf = nil
			}
			continue
		}
		if inCodeBlock {
			codeBuf = append(codeBuf, line)
			continue
		}

		trimmed := line
		// Heading
		if strings.HasPrefix(trimmed, "### ") {
			out = append(out, ansiBold+renderInline(trimmed[4:])+ansiReset)
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			out = append(out, ansiBold+ansiUnderline+renderInline(trimmed[3:])+ansiReset)
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			out = append(out, ansiBold+ansiUnderline+renderInline(trimmed[2:])+ansiReset)
			continue
		}

		// Horizontal rule
		if strings.TrimSpace(trimmed) == "---" {
			out = append(out, ansiFaint+strings.Repeat("─", width)+ansiReset)
			continue
		}

		// Unordered list
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			out = append(out, "  • "+renderInline(trimmed[2:]))
			continue
		}

		out = append(out, renderInline(trimmed))
	}
	return strings.Join(out, "\n")
}

func renderInline(s string) string {
	if s == "" {
		return ""
	}
	// Process inline formatting: code, bold, italic, links
	var buf strings.Builder
	i := 0
	for i < len(s) {
		// Inline code: `...`
		if s[i] == '`' {
			end := strings.Index(s[i+1:], "`")
			if end >= 0 {
				code := s[i+1 : i+1+end]
				buf.WriteString(ansiFaint)
				buf.WriteString(code)
				buf.WriteString(ansiReset)
				i += end + 2
				continue
			}
		}
		// Bold: **...**
		if i+1 < len(s) && s[i] == '*' && s[i+1] == '*' {
			end := findClosing(s, i+2, "**")
			if end >= 0 {
				inner := s[i+2 : end]
				buf.WriteString(ansiBold)
				buf.WriteString(renderInline(inner))
				buf.WriteString(ansiReset)
				i = end + 2
				continue
			}
		}
		// Italic: *...*
		if s[i] == '*' {
			end := findClosing(s, i+1, "*")
			if end >= 0 {
				inner := s[i+1 : end]
				buf.WriteString(ansiItalic)
				buf.WriteString(renderInline(inner))
				buf.WriteString(ansiReset)
				i = end + 1
				continue
			}
		}
		// Link: [text](url) - show text only
		if s[i] == '[' {
			endText := strings.Index(s[i:], "](")
			if endText >= 0 {
				endURL := strings.Index(s[i+endText+2:], ")")
				if endURL >= 0 {
					text := s[i+1 : i+endText]
					// url := s[i+endText+2 : i+endText+2+endURL]
					buf.WriteString(ansiUnderline)
					buf.WriteString(text)
					buf.WriteString(ansiReset)
					i += endText + 2 + endURL + 1
					continue
				}
			}
		}
		buf.WriteByte(s[i])
		i++
	}
	return buf.String()
}

func findClosing(s string, start int, delimiter string) int {
	for i := start; i < len(s); i++ {
		if strings.HasPrefix(s[i:], delimiter) {
			return i
		}
		// Skip escaped characters
		if s[i] == '\\' {
			i++
		}
	}
	return -1
}

func renderCodeBlock(lines []string, width int) string {
	if len(lines) == 0 {
		return ""
	}
	var buf strings.Builder
	// Top border
	buf.WriteString(ansiFaint + "┌" + strings.Repeat("─", width-2) + "┐" + ansiReset + "\n")
	for _, line := range lines {
		buf.WriteString(ansiFaint + "│ " + ansiReset)
		buf.WriteString(ansiItalic + line + ansiReset)
		buf.WriteString("\n")
	}
	// Bottom border
	buf.WriteString(ansiFaint + "└" + strings.Repeat("─", width-2) + "┘" + ansiReset)
	return buf.String()
}

// renderTranscriptMessage renders a single message for the transcript.
func renderTranscriptMessage(role, content string, width int) []string {
	lines := strings.Split(mdRender(content, width), "\n")
	if role == "user" {
		// User messages: prefix each line with "› "
		prefix := ansiBold + "› " + ansiReset
		result := make([]string, len(lines))
		for i, l := range lines {
			result[i] = prefix + l
		}
		return result
	}
	// Assistant messages: just the rendered markdown
	return lines
}
