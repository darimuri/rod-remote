package rodpipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rp "github.com/darimuri/rod-remote/rod_pipeline"
	"github.com/darimuri/rod-remote/rod_pipeline/task"
	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pipeline")
}

var (
	productUrl = "https://janginthe.com/product/%EC%9D%98%EC%A0%95%EB%B6%80-%EC%9E%A5%EC%9D%B8%ED%95%9C%EA%B3%BC-%EB%AA%BB%EB%82%9C%EC%9D%B4-%EC%95%BD%EA%B3%BC-%ED%8C%8C%EC%A7%80%EC%95%BD%EA%B3%BC/260/category/28/display/1/"
	//productUrl = "https://janginthe.com/product/%EC%9E%A5%EC%9D%B8%EB%8D%94-%EC%95%BD%EA%B3%BC%EB%B9%B5/258/category/24/display/1/"
)

func purchaseClickHandler(wait func() *proto.PageJavascriptDialogOpening, handle func(*proto.PageHandleJavaScriptDialog) error) {
	wait()

	var noDialog bool

	errHandle := handle(&proto.PageHandleJavaScriptDialog{Accept: false, PromptText: ""})
	if errHandle != nil {
		if errHandle == context.DeadlineExceeded {
		} else if cdpError, ok := errHandle.(*cdp.Error); ok {
			switch cdpError.Code {
			case -32602: //No dialog is showing
				noDialog = true
			case -32001: //Session with given id not found
			default:
				panic(cdpError)

			}
		} else if errHandle.Error() != "context canceled" {
			panic(errHandle)
		}
	}

	if noDialog {
		time.Sleep(time.Second)
	}
}

var _ = Describe("janginthe.com purchase", Ordered, func() {
	var b *rod.Browser
	var p *rod.Page
	var cut *rp.Pipeline

	var loginFormTasks []types.ITask
	var loginAllTasks []types.ITask

	var logoutCondition types.ConditionalFunc

	BeforeEach(func() {
		b = rod.New().ControlURL(rp.ControlUrl)
		err := b.Connect()
		Expect(err).NotTo(HaveOccurred())

		p, err = b.Page(proto.TargetCreateTarget{})
		Expect(err).NotTo(HaveOccurred())
		Expect(p).NotTo(BeNil())

		bounds := p.MustGetWindow()
		err = p.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: *bounds.Width, Height: *bounds.Height})
		Expect(err).NotTo(HaveOccurred())

		loginFormTasks = rp.Tasks(
			task.Input("#member_id", rp.TestId),
			task.Input("#member_passwd", rp.TestPass),
			task.Click("div.login > fieldset > a.btn_login", nil),
			task.WaitLoad(),
			task.WaitIdle(time.Minute),
		)

		loginAllTasks = rp.Tasks(
			task.Click("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogoff.right.off > li:nth-child(1) > a", nil),
			task.WaitLoad(),
			task.WaitIdle(time.Minute),
		)
		loginAllTasks = append(loginAllTasks, loginFormTasks...)

		logoutCondition = task.Has("#header > div.inner > div > div > ul.xans-element-.xans-layout.xans-layout-statelogoff.right.off")

		cut = rp.NewPipeline(p)
	})

	AfterEach(func() {
		err := p.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("purchase using credit card", func() {
		purchaseTasks := rp.Tasks(
			task.Click("#addr_paymethod1", nil),
			task.Click("#btn_payment", nil),
		)

		testCut(cut, logoutCondition, loginAllTasks, loginFormTasks, purchaseTasks)
	})

	It("purchase using banking", func() {
		purchaseTasks := rp.Tasks(
			task.Click("#addr_paymethod0", nil),
			task.Input("#pname", "정구현"),
			task.Click("#allAgree", nil),
			task.Click("#bankaccount", nil),
			task.Type(input.ArrowDown),
			task.Type(input.Enter),
			task.Click("#btn_payment", nil),
		)

		testCut(cut, logoutCondition, loginAllTasks, loginFormTasks, purchaseTasks)
	})
})

func testCut(cut *rp.Pipeline, logoutCondition types.ConditionalFunc, loginAllTasks, loginFormTasks, purchaseTasks []types.ITask) {
	cut.
		Open(productUrl).
		WaitLoad().
		WaitIdle(time.Minute).
		If(
			logoutCondition,
			loginAllTasks,
			rp.Else(),
		).
		If(task.Visible("div.ec-base-button > a.first"), rp.Then(), rp.Else(
			task.RemoveClass("div.ec-base-button > a.first", "displaynone"),
			task.RemoveClass("div.ec-base-button > a.btnWhite", "displaynone"),
			task.AddClass("div.ec-base-button > span.btnBlack", "displaynone"),
		)).
		Input("#quantity", "1").
		While(task.Visible("#frm_order_act"), rp.Then(), rp.Else(
			task.Click("div.ec-base-button > a.first", purchaseClickHandler),
			task.Custom(func(p *rod.Page) error {
				//time.Sleep(time.Millisecond * 100)
				return nil
			}),
			task.WaitLoad(),
			task.WaitIdle(time.Minute),
			task.If(logoutCondition, loginFormTasks, rp.Else()),
		), 1000000).
		If(
			task.Has("#frm_order_act"),
			purchaseTasks,
			rp.Else(task.Stop("order form not found by condition #frm_order_act")),
		)

	err := cut.Run()
	Expect(err).NotTo(HaveOccurred())
}
