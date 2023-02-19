package main

import (
	"github.com/CyrilSbrodov/metricService.git/internal/app/protoServer"
)

func main() {
	srv := protoServer.NewServerApp()
	srv.Run()
}
