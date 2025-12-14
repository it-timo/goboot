package goboottypes

// Registrar defines a service interface capable of receiving script-related registrations.
//
// It is typically implemented by script-generating services like `base_local`
// that aggregate lines of logic from other services and render them into structured files.
//
// Both `RegisterLines` and `RegisterFile` are expected to be idempotent and
// reject duplicate entries for the same name.
type Registrar interface {
	// RegisterLines registers a named group of commands or script lines
	// to be rendered into a multiservice target (e.g., Makefile, Taskfile).
	//
	// The `name` should reflect the originating service and must be unique
	// within the output format. Returns an error on conflict or validation failure.
	RegisterLines(name string, lines []string) error

	// RegisterFile registers a script file to be created with the given name and content lines.
	//
	// Unlike RegisterLines (which groups entries by service), this function allows
	// direct control over the file name and script content. Used for standalone scripts.
	RegisterFile(name string, lines []string) error
}

// ScriptReceiver defines an interface for services that accept a Registrar
// to delegate script registration during generation.
//
// This enables one service (e.g., a linter or formatter module) to contribute
// logic to another service (e.g., base_local) without hard-coding dependencies.
//
// ScriptReceiver is generally called during the setup or `RegisterServices` phase.
type ScriptReceiver interface {
	// SetScriptReceiver provides the implementing service with a Registrar instance,
	// enabling it to submit script content dynamically during generation.
	SetScriptReceiver(registrar Registrar)
}
