package goboottypes

// DirPerm is the default directory mode for generated directories.
const DirPerm = 0755

// FilePerm is the default mode for generated non-executable files.
const FilePerm = 0644

// ScriptPerm is the default mode for generated executable script files.
const ScriptPerm = 0755

// TemplateSuffix marks template input files before rendering.
const TemplateSuffix = ".tmpl"
