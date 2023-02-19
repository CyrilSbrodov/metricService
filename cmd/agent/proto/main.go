package main

import "github.com/CyrilSbrodov/metricService.git/internal/app/protoAgent"

func main() {
	c := protoAgent.NewAgentApp()
	c.Run()
}
