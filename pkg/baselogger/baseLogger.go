/*
Package baselogger validates logger settings for generated project code.
*/
package baselogger

import (
	"errors"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/rs/zerolog"
)

// BaseLogger validates logger settings. File-owning services render logger-aware files.
type BaseLogger struct {
	cfg *config.BaseLoggerConfig
	log zerolog.Logger
}

// NewBaseLogger constructs BaseLogger for a target directory.
func NewBaseLogger(_ string) *BaseLogger {
	return &BaseLogger{
		log: zerolog.Nop(),
	}
}

// SetLogger injects the service-specific logger.
func (b *BaseLogger) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// ID returns the service identifier.
func (b *BaseLogger) ID() string {
	return goboottypes.ServiceNameBaseLogger
}

// SetConfig assigns validated base_logger config.
func (b *BaseLogger) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseLoggerConfig)
	if !ok {
		return errors.New("invalid config type for base_logger")
	}

	b.cfg = baseCfg

	return nil
}

// Run records the selected logger settings. Rendering is owned by base_project.
func (b *BaseLogger) Run() error {
	if b.cfg == nil {
		return errors.New("base_logger config is not assigned")
	}

	b.log.Info().
		Str("project_name", b.cfg.ProjectName).
		Str("logger_type", b.cfg.LoggerType).
		Msg("base_logger settings selected")

	return nil
}
