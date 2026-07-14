/*
Package goboottypes defines shared constants and interfaces used across goboot.
*/
package goboottypes

// Service IDs used in config and orchestration.
const (
	// ServiceNameBaseProject identifies the base project generator.
	ServiceNameBaseProject = "base_project"
	// ServiceNameBaseLint identifies the lint generator.
	ServiceNameBaseLint = "base_lint"
	// ServiceNameBaseLocal identifies the local tooling generator.
	ServiceNameBaseLocal = "base_local"
	// ServiceNameBaseTest identifies the test scaffolding generator.
	ServiceNameBaseTest = "base_test"
	// ServiceNameBaseCI identifies the CI generator.
	ServiceNameBaseCI = "base_ci"
	// ServiceNameBaseLogger identifies the logger wiring generator.
	ServiceNameBaseLogger = "base_logger"
	// ServiceNameBaseDocker identifies the containerization generator.
	ServiceNameBaseDocker = "base_docker"
	// ServiceNameBaseRelease identifies the release automation generator.
	ServiceNameBaseRelease = "base_release"
	// ServiceNameBaseGovernance identifies the governance generator.
	ServiceNameBaseGovernance = "base_governance"
	// ServiceNameBaseSupplyChain identifies the supply-chain security generator.
	ServiceNameBaseSupplyChain = "base_supplychain"
)
