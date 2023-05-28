package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/darimuri/rod-remote/server"
)

var (
	chromeBin = ""
	dataDir   = ""
	headless  = true
	port      = 9222
)

func init() {
	envChromeBin := "LAUNCHER_CHROME_BIN"
	envDataDir := "LAUNCHER_DATA_DIR"
	envNoHeadless := "LAUNCHER_NO_HEADLESS"
	envPort := "LAUNCHER_PORT"

	envs := []string{envChromeBin, envDataDir, envNoHeadless, envPort}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	log.Println("check environment variables", envs)

	if os.Getenv(envChromeBin) != "" {
		chromeBin = os.Getenv(envChromeBin)
	}

	if os.Getenv(envDataDir) != "" {
		dataDir = os.Getenv(envDataDir)
	}

	if os.Getenv(envNoHeadless) == "true" {
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
	var s *server.ApiServer
	var err error

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("close browser by signal %v\n", sig)
		if s != nil {
			err = s.Stop()
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

	s = server.New(":8080", chromeBin, dataDir, port, headless)
	s.Start()

}
