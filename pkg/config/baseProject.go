package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/it-timo/goboot/pkg/types"
)

// BaseProjectConfig defines the metadata used by goboot to generate a new Go project.
//
// It injects values into templates (e.g., README, LICENSE, CI configs) and governs
// how project-specific identity and versioning are rendered.
type BaseProjectConfig struct {
	// ProjectURL is the full repository URL (e.g., "https://github.com/user/project").
	// Used in go.mod and README links.
	ProjectURL string `yaml:"project_url"`

	// ProjectName is the main identifier for the project (e.g., "goboot").
	// Used in CLI entry points, module path, and template files.
	ProjectName string `yaml:"project_name"`

	// CapsProjectName is the uppercase variant of ProjectName (e.g., "GOBOOT").
	// Used in headers, LICENSE, and NOTICE.
	CapsProjectName string `yaml:"caps_project_name"`

	// LowerProjectName is the lowercase variant (e.g., "goboot").
	// Used for safe filenames, Docker images, etc.
	LowerProjectName string `yaml:"lower_project_name"`

	// UsedGoVersion specifies the Go version to write into config files (e.g., "1.22.2").
	UsedGoVersion string `yaml:"used_go_version"`

	// UsedNodeVersion specifies the Node.js version for optional tooling (e.g., "20.11.1").
	UsedNodeVersion string `yaml:"used_node_version"`

	// CurrentYear is injected into LICENSE and NOTICE. If not set, it defaults to the current system year.
	CurrentYear int `yaml:"current_year"`

	// ReleaseCurrentWindow is the current roadmap target (e.g., "Q2 2025").
	ReleaseCurrentWindow string `yaml:"release_current_window"`

	// ReleaseUpcomingWindow defines the next milestone window (e.g., "Q4 2025").
	ReleaseUpcomingWindow string `yaml:"release_upcoming_window"`

	// ReleaseLongTerm defines a long-term year-based goal horizon (e.g., "2028").
	ReleaseLongTerm string `yaml:"release_long_term"`

	// The Author is the project creator/owner, injected into LICENSE and docs.
	Author string `yaml:"author"`

	// GitProvider determines how to render related templates and links.
	GitProvider string `yaml:"git_provider"`

	// GitHubUser is the GitHub username or org (used in badges and URLs).
	GitHubUser string `yaml:"github_user"`

	// GitLabUser is the GitLab username or group (used in badges and URLs).
	GitLabUser string `yaml:"gitlab_user"`
}

// newBaseProjectConfig returns a newly initialized, empty BaseProjectConfig.
func newBaseProjectConfig() *BaseProjectConfig {
	return &BaseProjectConfig{}
}

// ID returns a stable identifier for this config.
func (bp *BaseProjectConfig) ID() string {
	return types.ServiceNameBaseProject
}

// ReadConfig loads the base project configuration from the provided YAML file path.
// It overwrites the current config values with the file contents.
func (bp *BaseProjectConfig) ReadConfig(confPath string) error {
	return readYMLConfig(confPath, bp)
}

// Validate verifies the BaseProjectConfig for use in scaffolding.
//
// It returns an error if required values are missing/invalid, or calls fillNeededInfos.
func (bp *BaseProjectConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bp.ProjectURL) == "" {
		missing = append(missing, "project_url")
	}

	if strings.TrimSpace(bp.ProjectName) == "" {
		missing = append(missing, "project_name")
	}

	if strings.TrimSpace(bp.UsedGoVersion) == "" {
		missing = append(missing, "used_go_version")
	}

	if strings.TrimSpace(bp.UsedNodeVersion) == "" {
		missing = append(missing, "used_node_version")
	}

	if strings.TrimSpace(bp.ReleaseCurrentWindow) == "" {
		missing = append(missing, "release_current_window")
	}

	if strings.TrimSpace(bp.ReleaseUpcomingWindow) == "" {
		missing = append(missing, "release_upcoming_window")
	}

	if strings.TrimSpace(bp.ReleaseLongTerm) == "" {
		missing = append(missing, "release_long_term")
	}

	if strings.TrimSpace(bp.Author) == "" {
		missing = append(missing, "author")
	}

	if strings.TrimSpace(bp.GitProvider) == "github" {
		if strings.TrimSpace(bp.GitHubUser) == "" {
			missing = append(missing, "github_user")
		}
	} else if strings.TrimSpace(bp.GitProvider) == "gitlab" {
		if strings.TrimSpace(bp.GitLabUser) == "" {
			missing = append(missing, "gitlab_user")
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

	// Autofill year if not set.
	if bp.CurrentYear == 0 {
		bp.CurrentYear = time.Now().Year()
	}
}
