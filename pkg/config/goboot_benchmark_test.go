package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/it-timo/goboot/pkg/config"
)

func BenchmarkGoBootConfigParsing(benchmark *testing.B) {
	const serviceCount = 500

	configPath := filepath.Join(benchmark.TempDir(), "goboot.yml")
	serviceEntry := "  - id: disabled_service\n    enabled: false\n"
	content := "projectName: BenchmarkProject\n" +
		"targetPath: /tmp/goboot-benchmark\n" +
		"parallelism: 4\n" +
		"services:\n" +
		strings.Repeat(serviceEntry, serviceCount)

	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		benchmark.Fatal(err)
	}

	benchmark.ReportAllocs()
	benchmark.ReportMetric(serviceCount, "services/config")
	benchmark.ResetTimer()

	for benchmark.Loop() {
		rootConfig := config.NewGoBoot(configPath)
		if err := rootConfig.Init(); err != nil {
			benchmark.Fatal(err)
		}
	}
}
