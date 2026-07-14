package basegovernance_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBaseGovernance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Governance Suite")
}
