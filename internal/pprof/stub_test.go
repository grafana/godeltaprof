package pprof

import (
	"go/types"
	"strings"
	"testing"
)
import "golang.org/x/tools/go/packages"

func TestSignatureExpandFinalInlineFrame(t *testing.T) {
	checkSignature(t, "runtime/pprof",
		"runtime_expandFinalInlineFrame",
		"func runtime/pprof.runtime_expandFinalInlineFrame(stk []uintptr) []uintptr")
	checkSignature(t, "github.com/pyroscope-io/godeltaprof/internal/pprof",
		"runtime_expandFinalInlineFrame",
		"func github.com/pyroscope-io/godeltaprof/internal/pprof.runtime_expandFinalInlineFrame(stk []uintptr) []uintptr")
}

func TestSignatureCyclesPerSecond(t *testing.T) {
	checkSignature(t, "runtime/pprof",
		"runtime_cyclesPerSecond",
		"func runtime/pprof.runtime_cyclesPerSecond() int64")
	checkSignature(t, "github.com/pyroscope-io/godeltaprof/internal/pprof",
		"runtime_cyclesPerSecond",
		"func github.com/pyroscope-io/godeltaprof/internal/pprof.runtime_cyclesPerSecond() int64")
}

func checkSignature(t *testing.T, pkg string, name string, expectedSignature string) {
	cfg := &packages.Config{
		Mode:  packages.NeedImports | packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax,
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range pkgs {
		if strings.Contains(p.ID, ".test") {
			continue
		}
		f := p.Types.Scope().Lookup(name)
		if f != nil {
			ff := f.(*types.Func)
			if ff.String() != expectedSignature {
				t.Fatalf("expected %s, got %s", expectedSignature, ff.String())
			}
		}
	}

}
