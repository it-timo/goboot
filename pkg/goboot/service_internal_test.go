package goboot

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

type recordingService struct {
	id             string
	configSet      bool
	runCalled      bool
	runHook        func()
	setConfigHook  func(config.ServiceConfig)
	setConfigError error
	runError       error
}

// recordingScriptService extends recordingService to track script injection.
type recordingScriptService struct {
	recordingService
	registrarSet bool
}

func (m *recordingScriptService) SetScriptReceiver(_ goboottypes.Registrar) {
	m.registrarSet = true
}

func (m *recordingService) ID() string {
	return m.id
}

func (m *recordingService) SetConfig(cfg config.ServiceConfig) error {
	m.configSet = true
	if m.setConfigHook != nil {
		m.setConfigHook(cfg)
	}

	if m.setConfigError != nil {
		return m.setConfigError
	}

	return nil
}

func (m *recordingService) Run() error {
	m.runCalled = true
	if m.runHook != nil {
		m.runHook()
	}

	if m.runError != nil {
		return m.runError
	}

	return nil
}

type mockServiceConfig struct {
	id string
}

func (m *mockServiceConfig) ID() string {
	return m.id
}

func (m *mockServiceConfig) ReadConfig(_, _ string) error {
	return nil
}

func (m *mockServiceConfig) Validate() error {
	return nil
}

var _ = Describe("serviceManager internals", func() {
	var (
		cfgMgr      *config.Manager
		testManager *serviceManager
	)

	BeforeEach(func() {
		cfgMgr = config.NewConfigManager()
		testManager = newServiceManager(cfgMgr)
	})

	It("assigns configs and runs matching services", func() {
		Expect(cfgMgr.Register(&mockServiceConfig{id: "custom"})).To(Succeed())

		svc := &recordingService{id: "custom"}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.configSet).To(BeTrue())
		Expect(svc.runCalled).To(BeTrue())
	})

	It("runs prior, main, and subsequent services in order", func() {
		Expect(cfgMgr.Register(&mockServiceConfig{id: goboottypes.ServiceNameBaseProject})).To(Succeed())
		Expect(cfgMgr.Register(&mockServiceConfig{id: "custom"})).To(Succeed())
		Expect(cfgMgr.Register(&mockServiceConfig{id: goboottypes.ServiceNameBaseLocal})).To(Succeed())

		var order []string
		services := []*recordingService{
			{id: goboottypes.ServiceNameBaseProject},
			{id: "custom"},
			{id: goboottypes.ServiceNameBaseLocal},
		}

		for _, svc := range services {
			svc.runHook = func() {
				order = append(order, svc.id)
			}
			Expect(testManager.register(svc)).To(Succeed())
		}

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(order).To(Equal([]string{
			goboottypes.ServiceNameBaseProject,
			"custom",
			goboottypes.ServiceNameBaseLocal,
		}))
	})

	It("returns an error when assigning config fails", func() {
		Expect(cfgMgr.Register(&mockServiceConfig{id: "failing"})).To(Succeed())

		svc := &recordingService{
			id:             "failing",
			setConfigError: errors.New("boom"),
		}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to set config"))
		Expect(svc.runCalled).To(BeFalse())
	})

	It("skips services without configuration", func() {
		svc := &recordingService{id: "custom"}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.runCalled).To(BeFalse())
	})

	It("skips prior services that are not registered", func() {
		err := testManager.runPriorServices()
		Expect(err).NotTo(HaveOccurred())
	})

	It("skips prior services without configuration", func() {
		svc := &recordingService{id: goboottypes.ServiceNameBaseProject}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.runCalled).To(BeFalse())
	})

	It("skips subsequent services without configuration", func() {
		svc := &recordingService{id: goboottypes.ServiceNameBaseLocal}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.runCalled).To(BeFalse())
	})

	It("errors when registering a duplicate service", func() {
		svc := &recordingService{id: "dup"}
		Expect(testManager.register(svc)).To(Succeed())
		err := testManager.register(&recordingService{id: "dup"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("already registered"))
	})

	It("skips script injection when no registrar is registered", func() {
		Expect(cfgMgr.Register(&mockServiceConfig{id: "scripted"})).To(Succeed())

		svc := &recordingScriptService{
			recordingService: recordingService{id: "scripted"},
		}
		Expect(testManager.register(svc)).To(Succeed())

		err := testManager.runAll()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.configSet).To(BeTrue())
		Expect(svc.runCalled).To(BeTrue())
		Expect(svc.registrarSet).To(BeFalse())
	})
})
