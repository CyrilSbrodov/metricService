package main

import (
	_ "net/http/pprof"

	"github.com/CyrilSbrodov/metricService.git/internal/app"
)

func main() {
	srv := app.NewServerApp()
	srv.Run()
}
