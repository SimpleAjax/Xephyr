package backend_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestXephyrBackend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xephyr Backend Test Suite")
}
