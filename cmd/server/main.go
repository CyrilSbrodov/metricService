package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/CyrilSbrodov/metricService.git/internal/app"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\nbuildDate: %s\nbuildCommit: %s\n", buildVersion, buildDate, buildCommit)

	srv := app.NewServerApp()
	srv.Run()
}
