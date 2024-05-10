package ctxprogress

import (
	"context"
)

func Range[D any](ctx context.Context, vals []D, iter func(context.Context, int, D) bool) {
	for i, val := range vals {
		Set[Empty](ctx, i, len(vals))

		if !iter(Context[Empty](ctx, Empty{}), i, val) {
			return
		}
	}

	Set[Empty](ctx, len(vals), len(vals))
}
