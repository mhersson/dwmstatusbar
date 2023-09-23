package dwmstatusbar_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDwmstatusbar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dwmstatusbar Suite")
}
