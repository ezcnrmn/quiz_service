package main

import (
	"flag"
	"log"

	"github.com/ezcnrmn/quiz_service/internal/app/utils/consts"
	"github.com/ezcnrmn/quiz_service/internal/pkg/app"
)

func main() {
	addr := flag.String("addr", consts.DEFAULT_ADDRESS, "http address")
	flag.Parse()

	a, err := app.New(*addr)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
