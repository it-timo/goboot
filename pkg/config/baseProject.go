package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseProjectConfig contains template inputs for base project generation.
type BaseProjectConfig struct {
	// SourcePath points to project templates.
	SourcePath string `yaml:"sourcePath"`

	// ProjectURL is the repository URL.
	ProjectURL string `yaml:"-"`

	// RepoPath is the URL without scheme, used for imports and module paths.
	RepoPath string `yaml:"-"`

	// ProjectName is the main project identifier.
	ProjectName string `yaml:"-"`

	// CapsProjectName is the uppercase variant of ProjectName.
	CapsProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant of ProjectName.
	LowerProjectName string `yaml:"-"`

	// GitProvider selects provider-specific template behavior.
	GitProvider string `yaml:"-"`

	// Profile identifies the selected project baseline in generated documentation.
	Profile string `yaml:"-"`

	// UsedGoVersion is the Go version written into generated files.
	UsedGoVersion string `yaml:"usedGoVersion"`

	// UsedNodeVersion is the Node.js version for optional tooling.
	UsedNodeVersion string `yaml:"usedNodeVersion"`

	// CurrentYear is used in license and notice files.
	CurrentYear int `yaml:"currentYear"`

	// ReleaseCurrentWindow is the active roadmap window label.
	ReleaseCurrentWindow string `yaml:"releaseCurrentWindow"`

	// ReleaseUpcomingWindow is the next roadmap window label.
	ReleaseUpcomingWindow string `yaml:"releaseUpcomingWindow"`

	// ReleaseLongTerm is the long-term roadmap horizon label.
	ReleaseLongTerm string `yaml:"releaseLongTerm"`

	// Author is the project owner shown in generated docs.
	Author string `yaml:"author"`

	// GitUser is the provider account/org used in generated links.
	GitUser string `yaml:"gitUser"`

	// Logger contains optional logger settings supplied by base_logger.
	Logger goboottypes.LoggerSettings `yaml:"-"`
}

// SetProfile injects the validated template profile.
func (bp *BaseProjectConfig) SetProfile(profile string) {
	bp.Profile = profile
}

// newBaseProjectConfig creates a BaseProjectConfig with the given project name.
func newBaseProjectConfig(projectName string) *BaseProjectConfig {
	return &BaseProjectConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bp *BaseProjectConfig) ID() string {
	return goboottypes.ServiceNameBaseProject
}

// ReadConfig loads base_project config from confPath.
func (bp *BaseProjectConfig) ReadConfig(confPath string, repoURL string, gitProvider string) error {
	bp.ProjectURL = repoURL
	bp.GitProvider = gitProvider

	return readYMLConfig(confPath, bp)
}

// Validate checks required fields and fills derived values.
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

	err := validateProjectName(bp.ProjectName)
	if err != nil {
		return err
	}

	bp.fillNeededInfos()

	return nil
}

// fillNeededInfos populates derived fields used by templates.
func (bp *BaseProjectConfig) fillNeededInfos() {
	// Normalize name variants.
	bp.CapsProjectName = strings.ToUpper(bp.ProjectName)
	bp.LowerProjectName = strings.ToLower(bp.ProjectName)
	bp.RepoPath = strings.TrimPrefix(bp.ProjectURL, "https://")
	bp.RepoPath = strings.TrimPrefix(bp.RepoPath, "http://")

	// Default year when omitted.
	if bp.CurrentYear == 0 {
		bp.CurrentYear = time.Now().Year()
	}
}
