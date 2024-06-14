package ctxprogress

import (
	"context"
)

type RangeProgress struct {
	Current int
	Total   int
}

func (s RangeProgress) Progress() (int, int) { return s.Current, s.Total }

func Range[D any](parentCtx context.Context, vals []D, iter func(context.Context, int, D) bool) {
	ctx := New(parentCtx)
	for i, val := range vals {
		Set(ctx, RangeProgress{Current: i, Total: len(vals)})

		if !iter(ctx, i, val) {
			return
		}
	}

	Set(ctx, RangeProgress{Current: len(vals), Total: len(vals)})
}
