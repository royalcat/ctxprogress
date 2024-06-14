package ctxprogress_test

import (
	"context"
	"testing"

	"github.com/royalcat/ctxprogress"
	"github.com/shoenig/test/must"
)

func TestRange(t *testing.T) {
	ctx := context.Background()
	ctx = ctxprogress.New(ctx)
	i := 0
	ctxprogress.AddCallbackTree(ctx, func(p ctxprogress.ProgressTree) {
		current, total := p.Progress()
		must.GreaterEq(t, current, total)
		must.Eq(t, 3, total)
		must.Eq(t, i, current)
		i++
	})
	arr := []int{1, 2, 3}
	ctxprogress.Range(ctx, []int{1, 2, 3},
		func(context.Context, int, int) bool { return true },
	)
	must.Eq(t, len(arr)+1, i)
}

func TestNested(t *testing.T) {
	ctx := context.Background()

	arr1 := []int{0, 1, 2}
	arr2 := []int{10, 11, 12, 13}

	ctx = ctxprogress.New(ctx)

	var lastProgress float64

	ctxprogress.AddCallback(ctx, func(p ctxprogress.Progress) {
		current, total := p.Progress()
		must.LessEq(t, total, current)
		curProgress := float64(current) / float64(total)
		must.GreaterEq(t, lastProgress, curProgress)
		lastProgress = curProgress
	})

	ctxprogress.Range(ctx, arr1,
		func(ctx context.Context, i int, _ int) bool {
			j := 0
			ctxprogress.AddCallbackTree(ctx, func(p ctxprogress.ProgressTree) {
				if i < len(p.Children)-1 {
					return
				}
				must.Len(t, i+1, p.Children)
				current, total := p.Children[i].Progress()
				must.GreaterEq(t, current, total)
				must.Eq(t, j, current)
				must.Eq(t, len(arr2), total)
			})
			ctxprogress.Range(ctx, arr2,
				func(ctx context.Context, jj int, _ int) bool {
					j++
					return true
				},
			)

			return true
		},
	)
}
