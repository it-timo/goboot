package goboottypes_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Coverage is expected to show [no statements] because this package provides constants and interfaces only.
func TestTypes(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Types Suite")
}
