package main

import "github.com/CyrilSbrodov/metricService.git/internal/app"

func main() {
	agent := app.NewAgentApp()
	agent.Run()
}
