package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseProjectConfig defines the metadata used by goboot to generate a new Go project.
//
// It injects values into templates (e.g., README, LICENSE, CI configs) and governs
// how project-specific identity and versioning are rendered.
type BaseProjectConfig struct {
	// SourcePath is the path where the project will walk to get the template files.
	SourcePath string `yaml:"sourcePath"`

	// ProjectURL is the full repository URL (e.g., "https://github.com/user/project").
	// Used in README links etc.
	ProjectURL string `yaml:"-"`

	// RepoPath is the full repository path (e.g., "github.com/user/project").
	// Used in go.mod and imports.
	RepoPath string `yaml:"-"`

	// ProjectName is the main identifier for the project (e.g., "goboot").
	// Used in CLI entry points, module path, and template files.
	ProjectName string `yaml:"-"`

	// CapsProjectName is the uppercase variant of ProjectName (e.g., "GOBOOT").
	// Used in headers, LICENSE, and NOTICE.
	CapsProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant (e.g., "goboot").
	// Used for safe filenames, Docker images, etc.
	LowerProjectName string `yaml:"-"`

	// UsedGoVersion specifies the Go version to write into config files (e.g., "1.22.2").
	UsedGoVersion string `yaml:"usedGoVersion"`

	// UsedNodeVersion specifies the Node.js version for optional tooling (e.g., "20.11.1").
	UsedNodeVersion string `yaml:"usedNodeVersion"`

	// CurrentYear is injected into LICENSE and NOTICE. If not set, it defaults to the current system year.
	CurrentYear int `yaml:"currentYear"`

	// ReleaseCurrentWindow is the current roadmap target (e.g., "Q2 2025").
	ReleaseCurrentWindow string `yaml:"releaseCurrentWindow"`

	// ReleaseUpcomingWindow defines the next milestone window (e.g., "Q4 2025").
	ReleaseUpcomingWindow string `yaml:"releaseUpcomingWindow"`

	// ReleaseLongTerm defines a long-term year-based goal horizon (e.g., "2028").
	ReleaseLongTerm string `yaml:"releaseLongTerm"`

	// The Author is the project creator/owner, injected into LICENSE and docs.
	Author string `yaml:"author"`

	// GitProvider determines how to render related templates and links.
	GitProvider string `yaml:"gitProvider"`

	// GitUser is the GitHub username or org (used in badges and URLs).
	GitUser string `yaml:"gitUser"`
}

// newBaseProjectConfig returns a newly initialized BaseProjectConfig with the project name.
func newBaseProjectConfig(projectName string) *BaseProjectConfig {
	return &BaseProjectConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bp *BaseProjectConfig) ID() string {
	return goboottypes.ServiceNameBaseProject
}

// ReadConfig loads the base project configuration from the provided YAML file path.
// It overwrites the current config values with the file contents.
func (bp *BaseProjectConfig) ReadConfig(confPath string, repoURL string) error {
	bp.ProjectURL = repoURL

	return readYMLConfig(confPath, bp)
}

// Validate verifies the BaseProjectConfig for use in scaffolding.
//
// It returns an error if required values are missing/invalid, or calls fillNeededInfos.
//
//nolint:cyclop // flat validation logic preferred for clarity and extensibility.
func (bp *BaseProjectConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bp.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bp.ProjectURL) == "" {
		missing = append(missing, "projectUrl")
	}

	if strings.TrimSpace(bp.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(bp.UsedGoVersion) == "" {
		missing = append(missing, "usedGoVersion")
	}

	if strings.TrimSpace(bp.UsedNodeVersion) == "" {
		missing = append(missing, "usedNodeVersion")
	}

	if strings.TrimSpace(bp.ReleaseCurrentWindow) == "" {
		missing = append(missing, "releaseCurrentWindow")
	}

	if strings.TrimSpace(bp.ReleaseUpcomingWindow) == "" {
		missing = append(missing, "releaseUpcomingWindow")
	}

	if strings.TrimSpace(bp.ReleaseLongTerm) == "" {
		missing = append(missing, "releaseLongTerm")
	}

	if strings.TrimSpace(bp.Author) == "" {
		missing = append(missing, "author")
	}

	if strings.TrimSpace(bp.GitProvider) != "" {
		if strings.TrimSpace(bp.GitUser) == "" {
			missing = append(missing, "gitUser")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	bp.fillNeededInfos()

	return nil
}

// fillNeededInfos fills any derived fields in the config.
func (bp *BaseProjectConfig) fillNeededInfos() {
	// Normalize caps/lower.
	bp.CapsProjectName = strings.ToUpper(bp.ProjectName)
	bp.LowerProjectName = strings.ToLower(bp.ProjectName)
	bp.RepoPath = strings.TrimPrefix(bp.ProjectURL, "https://")
	bp.RepoPath = strings.TrimPrefix(bp.RepoPath, "http://")

	// Autofill year if not set.
	if bp.CurrentYear == 0 {
		bp.CurrentYear = time.Now().Year()
	}
}
