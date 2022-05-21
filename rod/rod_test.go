package rod

import (
	"context"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rod")
}

var _ = Describe("Rod", Ordered, func() {
	Context("Timeout", func() {
		It("WaitLoad after Timeout 1ms causes context.DeadlineExceeded error", func() {
			page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
			page = page.Timeout(time.Millisecond)
			err := page.WaitLoad()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.DeadlineExceeded))
		})

		It("WaitLoad after Timeout 1m causes no error", func() {
			page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
			page = page.Timeout(time.Minute)
			err := page.WaitLoad()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Launch", func() {
		It("Specify bin", func() {
			binPath := "/snap/bin/chromium"
			launcherURL := launcher.New().RemoteDebuggingPort(9222).Headless(true).Bin(binPath).MustLaunch()

			browser := rod.New().ControlURL(launcherURL)
			ExpectPage(browser)
		})

		It("LookPath", func() {
			binPath, found := launcher.LookPath()
			Expect(found).To(BeTrue())

			controlURL := launcher.New().Bin(binPath).MustLaunch()
			browser := rod.New().ControlURL(controlURL)
			ExpectPage(browser)
		})

		It("ResolveURL", func() {
			binPath, found := launcher.LookPath()
			Expect(found).To(BeTrue())

			launcherURL := launcher.New().RemoteDebuggingPort(9222).Headless(true).Bin(binPath).MustLaunch()

			browser := rod.New().ControlURL(launcherURL)
			ExpectPage(browser)
		})
	})
})

func ExpectPage(browser *rod.Browser) {
	err := browser.Connect()
	Expect(err).NotTo(HaveOccurred())

	page := browser.MustPage("https://www.wikipedia.org/")
	Expect(page).NotTo(BeNil())

	err = page.Timeout(time.Second * 3).WaitLoad()
	Expect(err).NotTo(HaveOccurred())

	err = page.WaitIdle(time.Second * 3)
	Expect(err).NotTo(HaveOccurred())
}
