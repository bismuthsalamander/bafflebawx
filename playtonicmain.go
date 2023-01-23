package main

import (
	"fmt"
	"os"

	"github.com/bismuthsalamander/bafflebawx/playtonic"
)

func main() {
	inceptorUrl := os.Getenv("INCEPTOR_URL")
	if len(inceptorUrl) == 0 {
		inceptorUrl = "http://inceptor:8080"
	}
	//todo: add configurable port number?
	conf := playtonic.PlaytonicConfig{inceptorUrl}
	s, err := playtonic.Server(conf)
	if err == nil {
		s.Run(":8080")
	}
}
