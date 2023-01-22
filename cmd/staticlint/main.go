package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// определяем map подключаемых правил

	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		mychecks = append(mychecks, v.Analyzer)

	}
	multichecker.Main(
		mychecks...,
	)
}
