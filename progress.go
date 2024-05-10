package ctxprogress

import (
	"context"
)

type Progress interface {
	Progress() (current, total int)
}

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

type Callback func(Progress)

func AddCallback(ctx context.Context, clb Callback) {
	entry := getNode(ctx)
	entry.clbs = append(entry.clbs, clb)
}

type CallbackTree func(ProgressTree)

func AddCallbackTree(ctx context.Context, clb CallbackTree) {
	entry := getNode(ctx)
	entry.clbsTrs = append(entry.clbsTrs, clb)
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

func GetTree(ctx context.Context) ProgressTree {
	ti := ctx.Value(entryCtxKey)
	if ti == nil {
		return ProgressTree{}
	}

	e, ok := ti.(*progressEntry)
	if !ok || e == nil {
		return ProgressTree{}
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
