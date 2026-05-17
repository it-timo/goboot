package config

import (
	"fmt"
	"strings"
)

func validateProjectName(projectName string) error {
	trimmed := strings.TrimSpace(projectName)
	if trimmed == "" {
		return nil
	}

	if trimmed != projectName {
		return fmt.Errorf("projectName %q must not contain leading or trailing whitespace", projectName)
	}

	if !isASCIIAlpha(trimmed[0]) {
		return fmt.Errorf("projectName %q must start with an ASCII letter", projectName)
	}

	for index := 1; index < len(trimmed); index++ {
		char := trimmed[index]
		if isASCIIAlpha(char) || isASCIIDigit(char) {
			continue
		}

		return fmt.Errorf("projectName %q must contain only ASCII letters and digits", projectName)
	}

	return nil
}

func isASCIIAlpha(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func isASCIIDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
