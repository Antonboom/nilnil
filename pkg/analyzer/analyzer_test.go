package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/nilnil/pkg/analyzer"
)

func TestNilNil(t *testing.T) {
	pkgs := []string{
		"examples",
		"strange",
	}
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), pkgs...)
}

func TestNilNil_Unsafe(t *testing.T) {
	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("checked-types", "uintptr,unsafeptr"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "unsafe")
}
