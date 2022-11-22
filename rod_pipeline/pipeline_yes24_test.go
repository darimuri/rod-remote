package rodpipeline_test

import (
	"fmt"
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
	"github.com/darimuri/rod-remote/rod_pipeline/userod"
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
			productId := 44166
			session := "18시 00분"
			sessionDate := "#\\32 023-02-04"

			url := fmt.Sprintf("http://m.ticket.yes24.com/Perf/Detail/PerfInfo.aspx?IdPerf=%d", productId)

			loginTasks := rp.Tasks(
				task.Tap("#wingScroll_wrap > div > div.greetingMsg > span.btn > a", nil),
				task.WaitLoad(),
				task.WaitIdle(time.Second*60),
				task.Input("#SMemberID", rp.TestId),
				task.Input("#SMemberPassword", rp.TestPass),
				task.Tap("#btn_login", nil),
			)

			reloadPageTasks := rp.Tasks(
				task.Reload(),
				task.WaitLoad(),
				task.WaitIdle(time.Minute),
				task.Custom(func(pc *types.PipelineContext) error {
					time.Sleep(time.Second)
					return nil
				}),
			)

			checkTimeFunc := func(_ *types.PipelineContext, el *rod.Element) (bool, error) {
				a, err := el.Attribute("timeinfo")
				if err != nil {
					return true, err
				}

				if a != nil && strings.TrimSpace(*a) == session {
					return true, el.Tap()
				}

				return false, nil
			}

			var reserved bool

			reserveButtonSelector := "#gd_norInfo > div.gd_btn > ul > li:nth-child(1) > a"
			//reserveButtonSelector := "#gd_norInfo > div.gd_btn > a.btn_c.btn_buy.btn_red"

			cut.
				Open(url).
				WaitLoad().
				WaitIdle(time.Minute).
				Tap("#entWing > span", nil).
				If(
					task.ContainsText("#wingScroll_wrap > div > div.greetingMsg > span.btn > a > em", "로그인"),
					loginTasks, rp.Else(task.Tap("#entWing > span", nil)),
				).
				While(
					task.Visible(reserveButtonSelector),
					rp.Then(task.Tap(reserveButtonSelector, nil)),
					reloadPageTasks,
					10000,
				).
				WaitLoad().
				WaitIdle(time.Minute).
				Tap(sessionDate, nil).
				WaitLoad().
				WaitIdle(time.Minute).
				ForEach("#ulTime > li", checkTimeFunc).
				While(
					func(pc *types.PipelineContext) (bool, error) {
						return reserved, nil
					},
					rp.Then(),
					rp.Else(
						task.If(
							task.Visible("#StepCtrlStep02_01 > div.guideTitArea > div.btnArea > a"),
							rp.Then(
								task.Custom(func(pc *types.PipelineContext) error {

									return nil
								}),
								task.Tap("#StepCtrlStep02_01 > div.guideTitArea > div.btnArea > a", nil),
								task.Tap("#grade_VIP석", nil),
							),
							rp.Else(),
						),
						task.ForEach("#seatSelDlScl > dl > dd:nth-child(2) > ul > li:nth-child(1) > a", func(_ *types.PipelineContext, el *rod.Element) (bool, error) {
							txt, err := el.Text()
							if err != nil {
								return true, err
							}
							s := strings.TrimSpace(txt)
							if strings.Contains(s, "석") {
								return true, el.Tap()
							}

							return false, nil
						}),
						task.Custom(func(pc *types.PipelineContext) error {
							el, err := pc.Page().Element("#ifrmSeatFrame")
							if err != nil {
								return err
							}
							pg, errFrame := el.Frame()
							if errFrame != nil {
								return errFrame
							}

							pc.PushPage(pg)

							return nil
						}),
						task.While(task.Has("#dMapInfo > map > area"),
							rp.Then(
								task.ForEach("#dMapInfo > map > area", func(pc *types.PipelineContext, el *rod.Element) (bool, error) {
									id, err := el.Attribute("id")
									if err != nil {
										return false, err
									} else if id == nil {
										return false, nil
									}

									img, err := pc.Page().Element("#blockFile")
									if err != nil {
										return true, err
									}
									center, err := userod.GetImageAreaCentroid(img, el)
									if err != nil {
										return true, err
									}

									errMoveTo := pc.Page().Mouse.MoveTo(*center)
									if errMoveTo != nil {
										return true, errMoveTo
									}

									//https://stackoverflow.com/questions/4529957/get-position-of-map-areahtml
									//get position of image and coords of area, then move mouse to the center of area

									//if *id == "area2" {
									//	return userod.EvalEventScript(el, "onclick")
									//}

									return false, nil
								}),
							),
							rp.Else(
								task.Custom(func(pc *types.PipelineContext) error {
									time.Sleep(time.Millisecond * 10)
									return nil
								}),
							), 1000),
						task.ForEach("#divSeatArray > div.s8", func(_ *types.PipelineContext, el *rod.Element) (bool, error) {
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

							_, err = el.Eval("() => ClickSeat(this)")
							if err != nil {
								return false, nil
							}

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

							if block != "4" && block != "5" {
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
						task.Custom(func(pc *types.PipelineContext) error {
							return pc.PopPage()
						}),
					),
					10000,
				)

			err := cut.Run()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
