package gobootutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// RenderTemplateToFile reads a file from the given fsRoot, renders it as a Go text/template,
// and writes the rendered content back to the same file location.
//
// Parameters:
//   - name: A template identifier (used for debugging and template naming).
//   - fsRoot: A secure *os.Root filesystem used for isolated, scoped access.
//   - path: The relative file path within fsRoot to be rendered.
//   - data: The data context for rendering (passed to template.Execute).
//
// Behavior:
//   - The function reads the file's content as raw text.
//   - Calls ExecuteTemplateText.
//   - Overwrites the file with the rendered result inside the fsRoot.
//
// Returns an error if the file cannot be read, parsed, rendered, or written.
//
// Notes:
//   - Only non-directory files should be passed.
//   - This method assumes the file exists before rendering.
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

// ExecuteTemplateText parses and renders a Go template from a raw string.
//
// Parameters:
//   - name: A template identifier (used for naming/debugging).
//   - text: The raw Go template source.
//   - data: The data context passed to template execution.
//
// Returns:
//   - The rendered string.
//   - An error if the template fails to parse or execute.
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

// templateFuncs returns the common function map used across goboot templates.
//
// Note: Only add small, deterministic helpers here to avoid surprising template behavior.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"indent":  indent,
		"oneLine": oneLine,
		"replace": strings.ReplaceAll,
	}
}

// indent prefixes every non-empty line in the provided string with the given number of spaces.
// It normalizes Windows line endings to Unix style before processing to keep behavior consistent across platforms.
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

// oneLine turns any newline / CRLF into single spaces, and trims the string.
func oneLine(curLine string) string {
	curLine = strings.ReplaceAll(curLine, "\r\n", "\n")
	curLine = strings.ReplaceAll(curLine, "\n", " ")

	return strings.TrimSpace(curLine)
}
