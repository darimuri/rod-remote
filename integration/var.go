package integration

import (
	"os"
)

var (
	NaverLogin    = ""
	NaverPassword = ""
)

func init() {
	if os.Getenv("NAVER_LOGIN") != "" {
		NaverLogin = os.Getenv("NAVER_LOGIN")
	}

	if os.Getenv("NAVER_PASSWORD") != "" {
		NaverPassword = os.Getenv("NAVER_PASSWORD")
	}
}
