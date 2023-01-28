package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/CyrilSbrodov/metricService.git/cmd/staticlint/checkosexit"
)

func main() {
	// создаем слайс анализаторов
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		mychecks = append(mychecks, v.Analyzer)
	}
	// определяем подключаемые анализаторы
	mychecks = append(mychecks,
		asmdecl.Analyzer, assign.Analyzer, atomic.Analyzer, atomicalign.Analyzer, bools.Analyzer,
		buildssa.Analyzer, buildtag.Analyzer, cgocall.Analyzer, composite.Analyzer, copylock.Analyzer,
		ctrlflow.Analyzer, deepequalerrors.Analyzer, errorsas.Analyzer, fieldalignment.Analyzer,
		findcall.Analyzer, framepointer.Analyzer, httpresponse.Analyzer, ifaceassert.Analyzer,
		loopclosure.Analyzer, lostcancel.Analyzer, nilfunc.Analyzer, nilness.Analyzer,
		pkgfact.Analyzer, printf.Analyzer, reflectvaluecompare.Analyzer, shadow.Analyzer,
		shift.Analyzer, sigchanyzer.Analyzer, sortslice.Analyzer, stdmethods.Analyzer,
		stringintconv.Analyzer, structtag.Analyzer, testinggoroutine.Analyzer, tests.Analyzer,
		timeformat.Analyzer, unmarshal.Analyzer, unreachable.Analyzer, unsafeptr.Analyzer,
		unusedresult.Analyzer, unusedwrite.Analyzer, usesgenerics.Analyzer)

	for _, ch := range stylecheck.Analyzers {
		mychecks = append(mychecks, ch.Analyzer)
	}
	// добавляем собственный анализатор
	mychecks = append(mychecks, checkosexit.Analyzer)

	multichecker.Main(
		mychecks...,
	)
}
