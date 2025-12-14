/*
Package goboottypes defines constants and foundational type declarations shared across the goboot project.

It provides a central place for declaring common identifiers such as service
names to ensure consistency across modules and reduce duplication.

This package is intentionally minimal and scoped to shared project-wide symbols.
*/
package goboottypes

// The declaration of service names.
const (
	// ServiceNameBaseProject is the name for the base project generation.
	ServiceNameBaseProject = "base_project"
	// ServiceNameBaseLint is the name for the base lint generation.
	ServiceNameBaseLint = "base_lint"
	// ServiceNameBaseLocal is the name for the base local generation.
	ServiceNameBaseLocal = "base_local"
	// ServiceNameBaseTest is the name for the base test generation.
	ServiceNameBaseTest = "base_test"
)
