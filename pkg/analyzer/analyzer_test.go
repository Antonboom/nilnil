//go:debug gotypesalias=1

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
	}
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), pkgs...)
}

func TestNilNil_Flags_CheckedTypes(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("checked-types", "ptr,uintptr"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "pointers-only")
}

func TestNilNil_Flags_DetectOpposite(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("detect-opposite", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "opposite")
}

func TestNilNil_Flags_DetectOppositeAndCheckedTypes(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("detect-opposite", "true"); err != nil {
		t.Fatal(err)
	}
	if err := anlzr.Flags.Set("checked-types", "chan,map"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "opposite-chan-map-only")
}

func TestNilNil_Flags_MultipleNils(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("only-two", "false"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "multiple-nils")
}

func TestNilNil_Flags_MultipleNilsAndDetectOpposite(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("only-two", "false"); err != nil {
		t.Fatal(err)
	}
	if err := anlzr.Flags.Set("detect-opposite", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), anlzr, "multiple-nils-and-opposite")
}
