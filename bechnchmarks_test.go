package ctxprogress_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/royalcat/ctxprogress"
)

func BenchmarkUpdateProgress(b *testing.B) {
	b.ReportAllocs()

	ctx := ctxprogress.New(context.Background())
	pprofFile, err := os.Create("mem.pprof")
	defer pprofFile.Close()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctxprogress.Set(ctx, ctxprogress.RangeProgress{Current: i, Total: b.N})
	}
	ctxprogress.Set(ctx, ctxprogress.RangeProgress{Current: b.N, Total: b.N})

	b.StopTimer()

	prg := ctxprogress.Get(ctx)
	current, total := prg.Progress()

	if current != total {
		b.Fatal(fmt.Errorf("invalid progress data, not completed"))
	}

	if total != b.N {
		b.Fatal(fmt.Errorf("invalid progress data, total not same as N"))
	}
}
