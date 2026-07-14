package regeneration_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRegeneration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Regeneration Suite")
}
