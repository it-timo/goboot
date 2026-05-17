package gobootutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// RenderTemplateToFile renders an existing file in fsRoot as a text/template and
// atomically replaces the same path with the rendered content.
func RenderTemplateToFile(name string, fsRoot *os.Root, path string, data any) error {
	raw, err := readRootFile(fsRoot, path)
	if err != nil {
		return err
	}

	rendered, err := ExecuteTemplateText(name, raw, data)
	if err != nil {
		return fmt.Errorf("failed template render: %w", err)
	}

	err = writeRootFileAtomic(fsRoot, path, rendered)
	if err != nil {
		return err
	}

	return nil
}

// WriteRootFile creates or replaces path inside fsRoot and returns only after
// write and close have both succeeded.
func WriteRootFile(fsRoot *os.Root, path string, content []byte, perm os.FileMode) error {
	outFile, err := fsRoot.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", path, err)
	}

	_, writeErr := outFile.Write(content)
	closeErr := outFile.Close()

	if writeErr != nil {
		return fmt.Errorf("failed to write file %q: %w", path, writeErr)
	}

	if closeErr != nil {
		return fmt.Errorf("failed to close file %q after write: %w", path, closeErr)
	}

	return nil
}

func readRootFile(fsRoot *os.Root, path string) (string, error) {
	file, err := fsRoot.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	raw, readErr := io.ReadAll(file)
	closeErr := file.Close()

	if readErr != nil {
		return "", fmt.Errorf("failed to read file: %w", readErr)
	}

	if closeErr != nil {
		return "", fmt.Errorf("failed to close file after read: %w", closeErr)
	}

	return string(raw), nil
}

func writeRootFileAtomic(fsRoot *os.Root, path string, content string) error {
	fileInfo, err := fsRoot.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat rendered target: %w", err)
	}

	tmpPath := filepath.Join(filepath.Dir(path), "."+filepath.Base(path)+".goboot-tmp")

	outFile, err := fsRoot.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed creation in root: %w", err)
	}

	_, writeErr := outFile.WriteString(content)
	closeErr := outFile.Close()

	if writeErr != nil {
		_ = fsRoot.Remove(tmpPath)

		return fmt.Errorf("failed to write rendered content: %w", writeErr)
	}

	if closeErr != nil {
		_ = fsRoot.Remove(tmpPath)

		return fmt.Errorf("failed to close rendered file: %w", closeErr)
	}

	err = fsRoot.Rename(tmpPath, path)
	if err != nil {
		_ = fsRoot.Remove(tmpPath)

		return fmt.Errorf("failed to replace rendered file: %w", err)
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
