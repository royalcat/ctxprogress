package ctxprogress

import (
	"context"
)

type Progress interface {
	Current() int
	Total() int
}

type Callback func(ProgressTree)

func nestProgress(ctx context.Context) context.Context {
	entry := &progressEntry{
		parent: getNode(ctx),
	}

	if entry.parent != nil {
		entry.parent.addChild(entry)
	}

	return context.WithValue(ctx, entryCtxKey, entry)
}

func Context(ctx context.Context) context.Context {
	return nestProgress(ctx)
}

func AddCallback(ctx context.Context, clb Callback) {
	entry := getNode(ctx)
	entry.clbs = append(entry.clbs, clb)
}

func Get(ctx context.Context) Progress {
	ti := ctx.Value(entryCtxKey)
	if ti == nil {
		return nil
	}

	e, ok := ti.(*progressEntry)
	if !ok || e == nil {
		return nil
	}

	return e.progress()
}

func Set(ctx context.Context, prg Progress) {
	entry := getNode(ctx)
	entry.current = prg
	entry.update()
}

func getNode(ctx context.Context) *progressEntry {
	ti := ctx.Value(entryCtxKey)
	if ti == nil {
		return nil
	}

	return ti.(*progressEntry)
}
