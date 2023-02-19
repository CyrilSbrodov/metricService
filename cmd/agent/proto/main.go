package main

import "github.com/CyrilSbrodov/metricService.git/internal/app/protoagent"

func main() {
	c := protoagent.NewAgentApp()
	c.Run()
}
