package rodpipeline_test

import (
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rp "github.com/darimuri/rod-remote/rod_pipeline"
	"github.com/darimuri/rod-remote/rod_pipeline/task"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pipeline")
}

var _ = Describe("janginthe.com purchase", Ordered, func() {
	var b *rod.Browser
	var p *rod.Page

	BeforeEach(func() {
		b = rod.New().ControlURL("ws://127.0.0.1:9222/devtools/browser/0dfdb7e6-9e36-44ab-a275-5903048e42e8")
		err := b.Connect()
		Expect(err).NotTo(HaveOccurred())

		p, err = b.Page(proto.TargetCreateTarget{})
		Expect(err).NotTo(HaveOccurred())

		bounds := p.MustGetWindow()
		err = p.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: *bounds.Width, Height: *bounds.Height})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := p.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("test1", func() {
		Expect(p).NotTo(BeNil())
		cut := rp.NewPipeline(p)

		loginFormTasks := task.Then(
			task.Input("#member_id", rp.TestId),
			task.Input("#member_passwd", rp.TestPass),
			task.Click("div.login > fieldset > a.btn_login"),
			task.WaitLoad(),
			task.WaitIdle(time.Minute),
		)

		loginAllTasks := task.Then(
			task.Click("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogoff.right.off > li:nth-child(1) > a"),
			task.WaitLoad(),
			task.WaitIdle(time.Minute),
		)
		loginAllTasks = append(loginAllTasks, loginFormTasks...)

		logoutCondition := task.Has("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogoff.right.off")

		productUrl := "https://janginthe.com/product/%EC%9D%98%EC%A0%95%EB%B6%80-%EC%9E%A5%EC%9D%B8%ED%95%9C%EA%B3%BC-%EB%AA%BB%EB%82%9C%EC%9D%B4-%EC%95%BD%EA%B3%BC-%ED%8C%8C%EC%A7%80%EC%95%BD%EA%B3%BC/260/category/28/display/1/"
		//productUrl := "https://janginthe.com/product/%EC%9E%A5%EC%9D%B8%EB%8D%94-%EC%95%BD%EA%B3%BC%EB%B9%B5/258/category/24/display/1/"

		cut.
			Open(productUrl).
			WaitLoad().
			WaitIdle(time.Minute).
			If(
				logoutCondition,
				loginAllTasks,
				task.Else(),
			).
			While(
				task.Visible("div.ec-base-button > a.first"),
				task.Then(),
				task.Else(
					task.Sleep(time.Second),
					task.Reload(),
					task.WaitLoad(),
					task.WaitIdle(time.Minute),
				),
				1000,
			).
			Click("div.ec-base-button > a.first").
			WaitLoad().
			WaitIdle(time.Minute).
			If(
				logoutCondition,
				loginFormTasks,
				task.Else(),
			).
			If(
				task.Has("#frm_order_act"),
				task.Then(
					task.Click("#addr_paymethod1"),
					task.Click("#btn_payment"),
				),
				task.Else(task.Stop("order form not found by condition #frm_order_act")),
			)

		err := cut.Run()
		Expect(err).NotTo(HaveOccurred())
	})
})
