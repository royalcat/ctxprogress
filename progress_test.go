package ctxprogress_test

import (
	"context"
	"testing"

	progress "github.com/royalcat/ctxprogress"
	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
	require := require.New(t)

	ctx := context.Background()
	ctx = progress.Context(ctx, progress.Empty{})
	i := 0
	progress.AddCallback(ctx, func(p progress.Progress[progress.Empty]) {
		require.LessOrEqual(p.Current, p.Total)
		require.Equal(3, p.Total)
		require.Equal(i, p.Current)
		i++
	})
	arr := []int{1, 2, 3}
	progress.Range(ctx, []int{1, 2, 3},
		func(context.Context, int, int) bool { return true },
	)
	require.Equal(len(arr)+1, i)
}

func TestNested(t *testing.T) {
	require := require.New(t)

	ctx := context.Background()

	ctx = progress.Context(ctx, progress.Empty{})
	i := 0
	progress.AddCallback(ctx, func(p progress.Progress[progress.Empty]) {
		require.LessOrEqual(p.Current, p.Total)
		require.Equal(3, p.Total)
		if p.Current == i {
			i++
		}

	})
	arr1 := []int{0, 1, 2}
	arr2 := []int{10, 11, 12, 13}

	progress.Range(ctx, arr1,
		func(ctx context.Context, _ int, _ int) bool {
			j := 0
			progress.AddCallback(ctx, func(p progress.Progress[progress.Empty]) {
				require.LessOrEqual(p.Current, p.Total)
				require.Equal(4, p.Total)
				require.Equal(j, p.Current)
				j++
			})
			progress.Range(ctx, arr2,
				func(ctx context.Context, _ int, _ int) bool {
					return true
				},
			)
			require.Equal(len(arr2)+1, j)
			return true
		},
	)
	require.Equal(len(arr1)+1, i)
}
