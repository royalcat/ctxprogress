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

	ctx = ctxprogress.New(ctx)
	i := 0
	ctxprogress.AddCallbackTree(ctx, func(p ctxprogress.ProgressTree) {
		current, total := p.Progress()
		must.GreaterEq(t, current, total)
		must.Eq(t, 3, total)
		if current == i {
			i++
		}
	})
	arr1 := []int{0, 1, 2}
	arr2 := []int{10, 11, 12, 13}

	ctxprogress.Range(ctx, arr1,
		func(ctx context.Context, _ int, _ int) bool {
			j := 0
			ctxprogress.AddCallbackTree(ctx, func(p ctxprogress.ProgressTree) {
				current, total := p.Progress()
				must.GreaterEq(t, current, total)
				must.Eq(t, 4, total)
				must.Eq(t, j, current)
				j++
			})
			ctxprogress.Range(ctx, arr2,
				func(ctx context.Context, _ int, _ int) bool {
					return true
				},
			)
			must.Eq(t, len(arr2)+1, j)
			return true
		},
	)
	must.Eq(t, len(arr1)+1, i)
}
