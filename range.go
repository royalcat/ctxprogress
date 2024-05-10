package ctxprogress

import (
	"context"
)

type RangeProgress struct {
	current int
	total   int
}

func (s RangeProgress) Current() int { return s.current }

func (s RangeProgress) Total() int { return s.total }

func Range[D any](ctx context.Context, vals []D, iter func(context.Context, int, D) bool) {
	for i, val := range vals {
		Set(ctx, RangeProgress{current: i, total: len(vals)})

		if !iter(Context(ctx), i, val) {
			return
		}
	}

	Set(ctx, RangeProgress{current: len(vals), total: len(vals)})
}
