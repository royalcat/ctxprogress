package ctxprogress

import (
	"context"
)

type Empty = struct{}

type Progress[P any] struct {
	Total      int
	Current    int
	Properties P

	Children []Progress[P]
}

type Callback[P any] func(Progress[P])

func nestProgress[P any](ctx context.Context, props P) context.Context {
	entry := &progressEntry[P]{
		properties: props,
		parent:     getEntry[P](ctx),
	}

	if entry.parent != nil {
		entry.parent.addChild(entry)
	}

	return context.WithValue(ctx, entryCtxKey, entry)
}

func Context[P any](ctx context.Context, props P) context.Context {
	return nestProgress[P](ctx, props)
}

func AddCallback[P any](ctx context.Context, clb Callback[P]) {
	entry := getEntry[P](ctx)
	entry.clbs = append(entry.clbs, clb)
}

func Get[P any](ctx context.Context) Progress[P] {
	ti := ctx.Value(entryCtxKey)
	if ti == nil {
		return Progress[P]{}
	}

	e, ok := ti.(*progressEntry[P])
	if !ok || e == nil {
		return Progress[P]{}
	}

	return e.progress()
}

func Set[P any](ctx context.Context, current, total int) {
	entry := getEntry[P](ctx)
	entry.current = current
	entry.total = total
	entry.update()
}

func SetProperties[P any](ctx context.Context, props P) {
	entry := getEntry[P](ctx)
	entry.properties = props
	entry.update()
}

func getEntry[P any](ctx context.Context) *progressEntry[P] {
	ti := ctx.Value(entryCtxKey)
	if ti == nil {
		return nil
	}

	return ti.(*progressEntry[P])
}
