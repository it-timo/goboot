package goboottypes

// Registrar collects generated command lines and script files.
type Registrar interface {
	// RegisterLines registers named command lines for shared outputs (for example Make/Task).
	RegisterLines(name string, lines []string) error

	// RegisterFile registers a script file by name and content lines.
	RegisterFile(name string, lines []string) error
}

// ScriptReceiver accepts a registrar used for local script outputs.
type ScriptReceiver interface {
	// SetScriptReceiver injects the registrar used for script registration.
	SetScriptReceiver(registrar Registrar)
}

// CIReceiver accepts a registrar used for CI job/file outputs.
type CIReceiver interface {
	// SetCIReceiver injects the registrar used for CI registration.
	SetCIReceiver(registrar Registrar)
}

// LoggerSettings describes the logger variant requested for generated project code.
type LoggerSettings struct {
	// Enabled reports whether generated runtime code should include logger wiring.
	Enabled bool

	// Type is the selected logger implementation.
	Type string
}

// LoggerSettingsProvider exposes logger settings to file-owning services.
type LoggerSettingsProvider interface {
	// LoggerSettings returns validated logger generation settings.
	LoggerSettings() LoggerSettings
}

// LoggerSettingsReceiver accepts logger settings before rendering owned files.
type LoggerSettingsReceiver interface {
	// SetLoggerSettings injects logger settings from base_logger.
	SetLoggerSettings(settings LoggerSettings)
}
