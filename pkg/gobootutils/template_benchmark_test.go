package gobootutils_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/it-timo/goboot/pkg/gobootutils"
)

const testProjectName = "goboot"

type templateBenchmarkData struct {
	Name    string
	Command string
}

func BenchmarkExecuteTemplateText(benchmark *testing.B) {
	testCases := []struct {
		name       string
		lineCount  int
		commandLen int
	}{
		{name: "small", lineCount: 8, commandLen: 32},
		{name: "large", lineCount: 512, commandLen: 512},
	}

	for _, testCase := range testCases {
		benchmark.Run(testCase.name, func(benchmark *testing.B) {
			rawTemplate := strings.Repeat("{{.Name}}: {{ oneLine .Command }}\n", testCase.lineCount)
			data := templateBenchmarkData{
				Name:    testProjectName,
				Command: strings.Repeat("x", testCase.commandLen),
			}

			benchmark.ReportAllocs()
			benchmark.ResetTimer()

			for benchmark.Loop() {
				rendered, err := gobootutils.ExecuteTemplateText("benchmark", rawTemplate, data)
				if err != nil {
					benchmark.Fatal(err)
				}

				runtime.KeepAlive(rendered)
			}
		})
	}
}

func BenchmarkRenderTemplateToFile(benchmark *testing.B) {
	root, err := os.OpenRoot(benchmark.TempDir())
	if err != nil {
		benchmark.Fatal(err)
	}

	benchmark.Cleanup(func() {
		if closeErr := root.Close(); closeErr != nil {
			benchmark.Error(closeErr)
		}
	})

	rawTemplate := strings.Repeat("project={{.Name}} command={{.Command}}\n", 128)
	data := templateBenchmarkData{Name: testProjectName, Command: "go test ./..."}

	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for benchmark.Loop() {
		if writeErr := gobootutils.WriteRootFile(root, "benchmark.txt", []byte(rawTemplate), 0o600); writeErr != nil {
			benchmark.Fatal(writeErr)
		}

		if renderErr := gobootutils.RenderTemplateToFile("benchmark", root, "benchmark.txt", data); renderErr != nil {
			benchmark.Fatal(renderErr)
		}
	}
}

func BenchmarkLargeProjectRendering(benchmark *testing.B) {
	const templateCount = 500

	rawTemplate := strings.Repeat("{{.Name}}: {{.Command}}\n", 64)
	data := templateBenchmarkData{Name: "service", Command: "go test -race ./..."}

	benchmark.ReportAllocs()
	benchmark.ReportMetric(templateCount, "templates/op")
	benchmark.ResetTimer()

	for benchmark.Loop() {
		generatedBytes := 0

		for templateIndex := range templateCount {
			rendered, err := gobootutils.ExecuteTemplateText("large-project", rawTemplate, data)
			if err != nil {
				benchmark.Fatalf("template %d: %v", templateIndex, err)
			}

			generatedBytes += len(rendered)
		}

		runtime.KeepAlive(generatedBytes)
	}
}
