package config_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

// Mock ServiceConfig for testing.
type mockServiceConfig struct {
	id          string
	shouldError bool
}

func (m *mockServiceConfig) ID() string {
	return m.id
}

func (m *mockServiceConfig) ReadConfig(_, _ string) error {
	if m.shouldError {
		return errors.New("mock read error")
	}

	return nil
}

func (m *mockServiceConfig) Validate() error {
	if m.shouldError {
		return errors.New("mock validation error")
	}

	return nil
}

var _ = Describe("Config Manager", func() {
	var manager *config.Manager

	BeforeEach(func() {
		manager = config.NewConfigManager()
	})

	Describe("NewConfigManager", func() {
		It("creates a new manager instance", func() {
			Expect(manager).NotTo(BeNil())
		})
	})

	Describe("Register", func() {
		Context("with a valid service config", func() {
			It("successfully registers the config", func() {
				mockCfg := &mockServiceConfig{
					id:          "test_service",
					shouldError: false,
				}

				err := manager.Register(mockCfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with an invalid service config (validation fails)", func() {
			It("returns an error and does not register", func() {
				mockCfg := &mockServiceConfig{
					id:          "invalid_service",
					shouldError: true,
				}

				err := manager.Register(mockCfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to validate config"))
			})
		})

		Context("when registering base_project service", func() {
			It("stores it as a registrar", func() {
				mockCfg := &mockServiceConfig{
					id:          goboottypes.ServiceNameBaseProject,
					shouldError: false,
				}

				err := manager.Register(mockCfg)
				Expect(err).NotTo(HaveOccurred())

				// Should be retrievable as registrar
				retrieved, ok := manager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeTrue())
				Expect(retrieved).To(Equal(mockCfg))

				_, ok = manager.GetService(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeFalse())
			})

			It("overwrites existing registrar entries on re-registration", func() {
				first := &mockServiceConfig{id: goboottypes.ServiceNameBaseProject}
				second := &mockServiceConfig{id: goboottypes.ServiceNameBaseProject}

				Expect(manager.Register(first)).To(Succeed())
				Expect(manager.Register(second)).To(Succeed())

				retrieved, ok := manager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeTrue())
				Expect(retrieved).To(Equal(second))
			})
		})

		Context("when registering other services", func() {
			It("stores them as regular services", func() {
				mockCfg := &mockServiceConfig{
					id:          "other_service",
					shouldError: false,
				}

				err := manager.Register(mockCfg)
				Expect(err).NotTo(HaveOccurred())

				// Should be retrievable as service
				retrieved, ok := manager.GetService("other_service")
				Expect(ok).To(BeTrue())
				Expect(retrieved).To(Equal(mockCfg))
			})
		})
	})

	Describe("GetService", func() {
		Context("when service exists", func() {
			It("returns the service and true", func() {
				mockCfg := &mockServiceConfig{
					id:          "existing_service",
					shouldError: false,
				}
				Expect(manager.Register(mockCfg)).To(Succeed())

				retrieved, ok := manager.GetService("existing_service")
				Expect(ok).To(BeTrue())
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.ID()).To(Equal("existing_service"))
			})
		})

		Context("when service does not exist", func() {
			It("returns nil and false", func() {
				retrieved, ok := manager.GetService("nonexistent")
				Expect(ok).To(BeFalse())
				Expect(retrieved).To(BeNil())
			})
		})
	})

	Describe("GetRegistrar", func() {
		Context("when registrar exists", func() {
			It("returns the registrar and true", func() {
				mockCfg := &mockServiceConfig{
					id:          goboottypes.ServiceNameBaseProject,
					shouldError: false,
				}
				Expect(manager.Register(mockCfg)).To(Succeed())

				retrieved, ok := manager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeTrue())
				Expect(retrieved).NotTo(BeNil())
			})
		})

		Context("when registrar does not exist", func() {
			It("returns nil and false", func() {
				retrieved, ok := manager.GetRegistrar("nonexistent")
				Expect(ok).To(BeFalse())
				Expect(retrieved).To(BeNil())
			})
		})
	})

	Describe("UnregisterService", func() {
		It("removes the service from the registry", func() {
			mockCfg := &mockServiceConfig{
				id:          "service_to_remove",
				shouldError: false,
			}
			Expect(manager.Register(mockCfg)).To(Succeed())

			// Verify it exists
			_, exist := manager.GetService("service_to_remove")
			Expect(exist).To(BeTrue())

			// Unregister it
			manager.UnregisterService("service_to_remove")

			// Verify it's gone
			_, exist = manager.GetService("service_to_remove")
			Expect(exist).To(BeFalse())
		})

		It("is a no-op if service doesn't exist", func() {
			// Should not panic or error
			Expect(func() {
				manager.UnregisterService("nonexistent")
			}).NotTo(Panic())
		})
	})

	Describe("UnregisterRegistrar", func() {
		It("removes the registrar from the registry", func() {
			mockCfg := &mockServiceConfig{
				id:          goboottypes.ServiceNameBaseProject,
				shouldError: false,
			}
			Expect(manager.Register(mockCfg)).To(Succeed())

			// Verify it exists
			_, exist := manager.GetRegistrar(goboottypes.ServiceNameBaseProject)
			Expect(exist).To(BeTrue())

			// Unregister it
			manager.UnregisterRegistrar(goboottypes.ServiceNameBaseProject)

			// Verify it's gone
			_, exist = manager.GetRegistrar(goboottypes.ServiceNameBaseProject)
			Expect(exist).To(BeFalse())
		})

		It("is a no-op if registrar doesn't exist", func() {
			// Should not panic or error
			Expect(func() {
				manager.UnregisterRegistrar("nonexistent")
			}).NotTo(Panic())
		})
	})

	Describe("ServiceConfigMeta", func() {
		Context("IsEnabled", func() {
			It("returns true when enabled is true", func() {
				meta := &config.ServiceConfigMeta{
					ID:       "test",
					ConfPath: "/path/to/config",
					Enabled:  true,
				}
				Expect(meta.IsEnabled()).To(BeTrue())
			})

			It("returns false when enabled is false", func() {
				meta := &config.ServiceConfigMeta{
					ID:       "test",
					ConfPath: "/path/to/config",
					Enabled:  false,
				}
				Expect(meta.IsEnabled()).To(BeFalse())
			})
		})
	})
})
