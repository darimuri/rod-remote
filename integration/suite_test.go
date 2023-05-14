package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/darimuri/rod-remote/control"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration test")
}

var testDevice = devices.Device{
	Title:          "Laptop with HiDPI screen",
	Capabilities:   []string{},
	UserAgent:      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36",
	AcceptLanguage: "ko,en;q=0.9,en-US;q=0.8",
	Screen: devices.Screen{
		DevicePixelRatio: 2,
		Horizontal: devices.ScreenSize{
			Width:  1920,
			Height: 1080,
		},
		Vertical: devices.ScreenSize{
			Width:  1080,
			Height: 1920,
		},
	},
}

var _ = Describe("nainom naver news profile reporter", Ordered, func() {
	var cut control.Control
	var profileURLs []string

	BeforeAll(func() {
		if NaverLogin == "" || NaverPassword == "" {
			Fail("NaverLogin, NaverPassword is required for this test")
		}

		launcherURL := "ws://127.0.0.1:9222/devtools/browser/059612b1-9c40-4de2-9c48-b8c3dc64dc30"

		browser := rod.New().ControlURL(launcherURL)
		err := browser.Connect()
		Expect(err).NotTo(HaveOccurred())

		pages, err := browser.Pages()
		Expect(err).NotTo(HaveOccurred())

		for _, p := range pages {
			p.MustClose()
		}

		cut = control.NewControl(browser)
		cut.DefaultDevice(testDevice.Landescape())
	})

	It("open nainom and get profile urls", func() {
		pc, err := cut.OpenPage("https://nainom.com/reports/naver_account", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(pc).NotTo(BeNil())

		pc.MustWindowMaximize()

		err = pc.Timeout(time.Second * 3).WaitLoad()
		Expect(err).NotTo(HaveOccurred())

		profileURLs, err = pc.GetAttributesFrom("ul > li > div > a", "href")
		Expect(err).NotTo(HaveOccurred())
		Expect(profileURLs).NotTo(HaveCap(0))

		pc.MustClose()
	})

	It("open naver and login", func() {
		pc, err := cut.OpenPage("https://naver.com", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(pc).NotTo(BeNil())

		pc.Timeout(time.Second * 5).MustWaitLoad()

		Expect(pc.MustHas("div#account")).To(BeTrue())
		accountDiv := pc.MustElement("div#account")
		Expect(accountDiv).NotTo(BeNil())
		classAttr := accountDiv.MustAttribute("class")
		Expect(classAttr).NotTo(BeNil())

		if strings.Contains(*classAttr, "sc_login") {
			Expect(*classAttr).To(ContainSubstring("sc_login")) //should be sc_my after login

			Expect(accountDiv.MustHas("a.link_login")).To(BeTrue())
			loginLink := accountDiv.MustElement("a.link_login")
			Expect(loginLink).NotTo(BeNil())

			clickedLink := loginLink.MustClick()
			Expect(clickedLink).NotTo(BeNil())

			pc.Timeout(time.Second * 5).MustWaitLoad()

			Expect(pc.MustHas("form#frmNIDLogin")).To(BeTrue())
			loginForm := pc.MustElement("form#frmNIDLogin")
			Expect(loginForm).NotTo(BeNil())

			Expect(loginForm.MustHas("input#id")).To(BeTrue())
			idInput := loginForm.MustElement("input#id")
			idInput.MustInput(NaverLogin)

			Expect(loginForm.MustHas("input#pw")).To(BeTrue())
			pwInput := loginForm.MustElement("input#pw")
			pwInput.MustInput(NaverPassword)

			Expect(loginForm.MustHas("label.keep_text")).To(BeTrue())
			keepLabel := loginForm.MustElement("label.keep_text")
			keepLabel.MustClick()

			Expect(loginForm.MustHas("button.btn_login")).To(BeTrue())
			loginButton := loginForm.MustElement("button.btn_login")
			loginButton.MustClick()

			pc.Timeout(time.Second * 5).MustWaitLoad()
		}

		Expect(pc.MustHas("div#account")).To(BeTrue())
		accountDiv = pc.MustElement("div#account")
		Expect(accountDiv).NotTo(BeNil())
		classAttr = accountDiv.MustAttribute("class")
		Expect(classAttr).NotTo(BeNil())
		Expect(*classAttr).To(ContainSubstring("sc_my"))

		pc.MustClose()
	})
})
