package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
)

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

var (
	chromeBin = ""
	dataDir   = ""
	headless  = true
	port      = 9222
)

func init() {
	envChromeBin := "LAUNCHER_CHROMEBIN"
	envDataDir := "LAUNCHER_DATADIR"
	envHeadless := "LAUNCHER_NOHEADLESS"
	envPort := "LAUNCHER_PORT"

	envs := []string{envChromeBin, envDataDir, envHeadless, envPort}

	log.Println("check environment variables", envs)

	if os.Getenv(envChromeBin) != "" {
		chromeBin = os.Getenv(envChromeBin)
	}

	if os.Getenv(envDataDir) != "" {
		dataDir = os.Getenv(envDataDir)
	}

	if os.Getenv(envHeadless) != "" {
		headless = false
	}

	if os.Getenv(envPort) != "" {
		parsed, err := strconv.Atoi(os.Getenv(envPort))
		if err != nil {
			panic(err)
		}
		port = parsed
	}
}

func main() {
	var browser *rod.Browser
	var err error

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("close browser by signal %v\n", sig)
		if browser != nil {
			err = browser.Close()
			if err != nil {
				panic(err)
			}
		}
		switch sig {
		case os.Interrupt:
			os.Exit(int(syscall.SIGINT))
		case syscall.SIGTERM:
			os.Exit(int(syscall.SIGTERM))
		default:
			log.Printf("unhandled signal case %v\n", sig)
			os.Exit(0)
		}
	}()

	if chromeBin == "" {
		binPath, found := launcher.LookPath()

		if found {
			log.Printf("found browser bin from %s\n", binPath)
			chromeBin = binPath
		} else {
			log.Printf("failed to find browser path\n")
		}

	}

	log.Println("launcher settings")
	log.Println("----------")
	log.Println("Bin :", chromeBin)
	log.Println("UserDataDir :", dataDir)
	log.Println("RemoteDebuggingPort :", port)
	log.Println("Headless :", headless)
	log.Println("----------")

	controlURL := launcher.New().RemoteDebuggingPort(port).Headless(headless).Bin(chromeBin).UserDataDir(dataDir).MustLaunch()

	log.Printf("launched browser with control url %s\n", controlURL)

	browser = rod.New().ControlURL(controlURL)
	err = browser.Connect()

	if err != nil {
		log.Panicf("failed to connect browser control url %s\n", controlURL)
	}

	select {}
}
