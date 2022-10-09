package control

import (
	"testing"

	"github.com/go-rod/rod"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "control")
}

var _ = Describe("XXX", Ordered, func() {
	var c Control
	var pc *PageControl

	BeforeEach(func() {
		b := rod.New().ControlURL("ws://127.0.0.1:9222/devtools/browser/66ce0aa2-2934-42c0-a668-f3c8937c41b0")
		err := b.Connect()
		Expect(err).NotTo(HaveOccurred())

		c = NewControl(b)
		pc, err = c.OpenPage("https://www.naver.com", true)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := pc.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("test1", func() {
		Expect(pc).NotTo(BeNil())
	})
})
