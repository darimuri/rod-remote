package rodpipeline_test

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rp "github.com/darimuri/rod-remote/rod_pipeline"
	"github.com/darimuri/rod-remote/rod_pipeline/task"
	"github.com/darimuri/rod-remote/rod_pipeline/types"
)

var _ = Describe("yes24", func() {
	Context("concert", func() {
		var b *rod.Browser
		var p *rod.Page
		var cut *rp.Pipeline

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

			cut = rp.NewPipeline(p)
		})

		AfterEach(func() {
			err := p.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		FIt("pyongchang", func() {
			url := "http://m.ticket.yes24.com/Perf/Detail/PerfInfo.aspx?IdPerf=43740"

			var reserved bool

			checkTimeFunc := func(el *rod.Element) (bool, error) {
				a, err := el.Attribute("timeinfo")
				if err != nil {
					return true, err
				}

				if a != nil && strings.TrimSpace(*a) == "18시 00분" {
					return true, el.Tap()
				}

				return false, nil
			}

			cut.
				Open(url).
				WaitLoad().
				WaitIdle(time.Minute).
				Tap("#entWing > span", nil).
				If(
					task.ContainsText("#wingScroll_wrap > div > div.greetingMsg > span.btn > a > em", "로그인"),
					rp.Then(
						task.Tap("#wingScroll_wrap > div > div.greetingMsg > span.btn > a", nil),
						task.WaitLoad(),
						task.WaitIdle(time.Second*60),
						task.Input("#SMemberID", rp.TestId),
						task.Input("#SMemberPassword", rp.TestPass),
						task.Tap("#btn_login", nil),
					), rp.Else(
						task.Tap("#entWing > span", nil),
					),
				).
				Tap("#gd_norInfo > div.gd_btn > a.btn_c.btn_buy.btn_red", nil).
				WaitLoad().
				WaitIdle(time.Minute).
				Tap("#\\32 022-12-31", nil).
				WaitLoad().
				WaitIdle(time.Minute).
				ForEachElement("#ulTime > li", checkTimeFunc).
				While(
					func(ctx *types.PipelineContext) (bool, error) {
						return reserved, nil
					},
					rp.Then(),
					rp.Else(
						task.Tap("#StepCtrlStep02_01 > div.guideTitArea > div.btnArea > a", nil),
						task.Tap("#grade_VIP석", nil),
						task.ForEach("#seatSelDlScl > dl > dd:nth-child(2) > ul > li:nth-child(1) > a", func(el *rod.Element) (bool, error) {
							txt, err := el.Text()
							if err != nil {
								return true, err
							}
							s := strings.TrimSpace(txt)
							if strings.Contains(s, "석") {
								return true, el.Click(proto.InputMouseButtonLeft, 1)
							}

							return false, nil
						}),
						task.Custom(func(p *rod.Page) error {
							el, err := p.Element("#ifrmSeatFrame")
							if err != nil {
								return err
							}
							pg, errFrame := el.Frame()
							if errFrame != nil {
								return errFrame
							}

							cut.PushPage(pg)

							return nil
						}),
						task.While(task.Has("#divSeatArray > div"), rp.Then(), rp.Else(task.Custom(func(p *rod.Page) error {
							time.Sleep(time.Millisecond * 10)
							return nil
						})), 1000),
						task.ForEach("#divSeatArray > div.s9", func(el *rod.Element) (bool, error) {
							a, err := el.Attribute("title")
							if err != nil {
								return true, err
							}
							if a == nil {
								return false, nil
							}
							s := strings.TrimSpace(*a)
							frc := strings.Split(s, " ")
							if len(frc) < 3 {
								return false, nil
							}
							floor := frc[0]
							br := frc[1]
							col := strings.ReplaceAll(frc[2], "번", "")

							if floor != "1층" {
								//TODO: go back to area selector
								return true, nil
							}

							blockAndRow := strings.Split(br, "구역")
							if len(blockAndRow) < 2 {
								return false, nil
							}

							block := blockAndRow[0]
							row := strings.ReplaceAll(blockAndRow[1], "열", "")

							if block == "가" {
								return true, nil
							}

							rowNum, errPRow := strconv.ParseInt(row, 10, 64)
							if errPRow != nil {
								return true, errPRow
							} else if rowNum >= 24 {
								//TODO: go back to area selector
								return true, nil
							}

							colNum, errPCol := strconv.ParseInt(col, 10, 64)
							if errPCol != nil {
								return true, errPCol
							} else if colNum > 3 {
								//try next seat
								return false, nil
							}

							if errTap := el.Tap(); errTap != nil {
								return true, errTap
							}
							reserved = true

							return true, nil
						}),
						task.Custom(func(p *rod.Page) error {
							return nil
						}),
					),
					10000,
				)

			err := cut.Run()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
