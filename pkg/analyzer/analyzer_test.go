package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/nilnil/pkg/analyzer"
)

func TestNilNil(t *testing.T) {
	t.Parallel()

	pkgs := []string{
		"examples",
		"strange",
		"unsafe",
	}
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), pkgs...)
}

func TestNilNil_Flags(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("checked-types", "ptr"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "pointers-only")
}
