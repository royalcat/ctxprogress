package ctxprogress_test

import (
	"context"
	"testing"

	"github.com/royalcat/ctxprogress"
	"github.com/shoenig/test/must"
)

func TestRange(t *testing.T) {
	ctx := context.Background()
	ctx = ctxprogress.Context(ctx)
	i := 0
	ctxprogress.AddCallback(ctx, func(p ctxprogress.ProgressTree) {
		must.GreaterEq(t, p.Current(), p.Total())
		must.Eq(t, 3, p.Total())
		must.Eq(t, i, p.Current())
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

	ctx = ctxprogress.Context(ctx)
	i := 0
	ctxprogress.AddCallback(ctx, func(p ctxprogress.ProgressTree) {
		must.GreaterEq(t, p.Current(), p.Total())
		must.Eq(t, 3, p.Total())
		if p.Current() == i {
			i++
		}
	})
	arr1 := []int{0, 1, 2}
	arr2 := []int{10, 11, 12, 13}

	ctxprogress.Range(ctx, arr1,
		func(ctx context.Context, _ int, _ int) bool {
			j := 0
			ctxprogress.AddCallback(ctx, func(p ctxprogress.ProgressTree) {
				must.GreaterEq(t, p.Current(), p.Total())
				must.Eq(t, 4, p.Total())
				must.Eq(t, j, p.Current())
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
