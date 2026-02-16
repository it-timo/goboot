package gobootutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// RenderTemplateToFile renders an existing file in fsRoot as a text/template and
// writes the rendered content back to the same path.
func RenderTemplateToFile(name string, fsRoot *os.Root, path string, data any) error {
	file, err := fsRoot.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer CloseFileWithErr(file)

	raw, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	rendered, err := ExecuteTemplateText(name, string(raw), data)
	if err != nil {
		return fmt.Errorf("failed template render: %w", err)
	}

	outFile, err := fsRoot.Create(path)
	if err != nil {
		return fmt.Errorf("failed creation in root: %w", err)
	}
	defer CloseFileWithErr(outFile)

	_, err = outFile.WriteString(rendered)
	if err != nil {
		return fmt.Errorf("failed to write rendered content: %w", err)
	}

	return nil
}

// ExecuteTemplateText parses and executes a template from raw text.
func ExecuteTemplateText(name, text string, data any) (string, error) {
	tmpl, err := template.New(name).Funcs(templateFuncs()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed template parse: %w", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed template execution: %w", err)
	}

	return buf.String(), nil
}

// templateFuncs returns deterministic helpers shared by template rendering.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"indent":  indent,
		"oneLine": oneLine,
		"replace": strings.ReplaceAll,
	}
}

// indent prefixes each non-empty line with spaces after normalizing CRLF to LF.
func indent(spaces int, curLine string) string {
	if spaces < 0 {
		spaces = 0
	}

	pad := strings.Repeat(" ", spaces)
	curLine = strings.ReplaceAll(curLine, "\r\n", "\n")
	lines := strings.Split(curLine, "\n")

	for index := 0; index < len(lines); index++ {
		if lines[index] == "" {
			continue
		}

		lines[index] = pad + lines[index]
	}

	return strings.Join(lines, "\n")
}

// oneLine replaces line breaks with spaces and trims surrounding whitespace.
func oneLine(curLine string) string {
	curLine = strings.ReplaceAll(curLine, "\r\n", "\n")
	curLine = strings.ReplaceAll(curLine, "\n", " ")

	return strings.TrimSpace(curLine)
}
