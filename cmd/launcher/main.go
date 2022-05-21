package main

import (
	"log"
	"os"
	"os/signal"
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

	binPath, found := launcher.LookPath()

	if found {
		log.Printf("found browser bin from %s\n", binPath)
	} else {
		log.Printf("failed to find browser path\n")
	}

	controlURL := launcher.New().RemoteDebuggingPort(9222).Headless(false).Bin(binPath).MustLaunch()

	log.Printf("launched browser with control url %s\n", controlURL)

	browser = rod.New().ControlURL(controlURL)
	err = browser.Connect()

	if err != nil {
		log.Panicf("failed to connect browser control url %s\n", controlURL)
	}

	select {}
}
