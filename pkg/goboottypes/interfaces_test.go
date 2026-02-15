package goboottypes_test

import (
	"github.com/it-timo/goboot/pkg/baseci"
	"github.com/it-timo/goboot/pkg/baselint"
	"github.com/it-timo/goboot/pkg/baselocal"
	"github.com/it-timo/goboot/pkg/basetest"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

// Compile-time interface assertions to guard regressions.
var (
	_ goboottypes.Registrar = (*baselocal.BaseLocal)(nil)
	_ goboottypes.Registrar = (*baseci.BaseCI)(nil)

	_ goboottypes.ScriptReceiver = (*baselint.BaseLint)(nil)
	_ goboottypes.ScriptReceiver = (*basetest.BaseTest)(nil)

	_ goboottypes.CIReceiver = (*baselint.BaseLint)(nil)
	_ goboottypes.CIReceiver = (*basetest.BaseTest)(nil)
)
