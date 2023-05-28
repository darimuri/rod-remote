package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	echopprof "github.com/hiko1129/echo-pprof"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// refer log format https://echo.labstack.com/middleware/logger/#configuration
const accessLogFormat = `${time_custom}: A ` +
	`${remote_ip}	${method}	${uri}	${status}	` +
	`${latency}	${latency_human}	${bytes_in}	${bytes_out}	` +
	`"${host}"	"${user_agent}"	"${error}"` + "\n"

const logTimeFormat = "2006-01-02 15:04:05.000000"

type ApiServer struct {
	listenAddress string
	chromeBin     string
	dataDir       string
	port          int
	headless      bool

	controlURL string
}

func New(listenAddress, chromeBin, dataDir string, port int, headless bool) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
		chromeBin:     chromeBin,
		dataDir:       dataDir,
		port:          port,
		headless:      headless,
	}
}

type stdWriter struct {
}

func (_ stdWriter) Write(p []byte) (n int, err error) {
	log.Println("rod:", strings.TrimSpace(string(p)))
	return len(p), nil
}

func (s *ApiServer) Start() {
	if s.chromeBin == "" {
		binPath, found := launcher.LookPath()

		if found {
			log.Printf("found browser bin from %s\n", binPath)
			s.chromeBin = binPath
		} else {
			log.Printf("failed to find browser path\n")
		}
	}

	log.Println("launcher settings")
	log.Println("----------")
	log.Println("Bin :", s.chromeBin)
	log.Println("UserDataDir :", s.dataDir)
	log.Println("RemoteDebuggingPort :", s.port)
	log.Println("Headless :", s.headless)
	log.Println("----------")

	s.controlURL = launcher.New().Logger(stdWriter{}).RemoteDebuggingPort(s.port).
		Set("enable-automation", "false").
		Set("no-first-run").
		Set("password-store", "basic").
		Set("use-mock-keychain").
		Set("start-maximized").
		Headless(s.headless).Bin(s.chromeBin).UserDataDir(s.dataDir).MustLaunch()

	//log.Printf("launched browser with control url %s\n", controlURL)

	e := echo.New()
	e.HideBanner = true

	e.Use(echoprometheus.NewMiddleware("echo"))

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           accessLogFormat,
		CustomTimeFormat: logTimeFormat,
	}))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/v1/control/url", func(c echo.Context) error {
		return c.String(http.StatusOK, s.controlURL)
	})

	echopprof.Wrap(e)

	// Start server
	e.Logger.Fatal(e.Start(s.listenAddress))
}

func (s *ApiServer) Stop() error {
	if s.controlURL != "" {
		b := rod.New().ControlURL(s.controlURL)
		err := b.Connect()
		if err != nil {
			return err
		}

		return b.Close()
	}

	return nil
}
