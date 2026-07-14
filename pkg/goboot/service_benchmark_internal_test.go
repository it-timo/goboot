package goboot

import (
	"crypto/sha256"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/rs/zerolog"
)

func BenchmarkRegularServiceExecution(benchmark *testing.B) {
	const serviceCount = 64

	for _, parallelism := range []int{1, 4, 8} {
		benchmark.Run(fmt.Sprintf("parallelism_%d", parallelism), func(benchmark *testing.B) {
			configManager := config.NewConfigManager()
			manager := newServiceManager(configManager, zerolog.Nop(), parallelism)
			payload := []byte(strings.Repeat("goboot-generation-workload", 256))

			for serviceIndex := range serviceCount {
				serviceID := fmt.Sprintf("benchmark-service-%03d", serviceIndex)
				if err := configManager.Register(&mockServiceConfig{id: serviceID}); err != nil {
					benchmark.Fatal(err)
				}

				service := &recordingService{
					id: serviceID,
					runHook: func() {
						digest := sha256.Sum256(payload)
						runtime.KeepAlive(digest)
					},
				}

				if err := manager.register(service); err != nil {
					benchmark.Fatal(err)
				}
			}

			if err := manager.assignConfigs(); err != nil {
				benchmark.Fatal(err)
			}

			benchmark.ReportAllocs()
			benchmark.ReportMetric(serviceCount, "services/op")
			benchmark.ResetTimer()

			for benchmark.Loop() {
				if err := manager.runRegularServices(); err != nil {
					benchmark.Fatal(err)
				}
			}
		})
	}
}
