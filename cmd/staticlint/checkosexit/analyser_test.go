package checkosexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyser(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
