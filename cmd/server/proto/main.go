package main

import (
	"github.com/CyrilSbrodov/metricService.git/internal/app/protoserver"
)

func main() {
	srv := protoserver.NewServerApp()
	srv.Run()
}
