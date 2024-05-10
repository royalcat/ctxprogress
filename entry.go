package ctxprogress

type progressKey uint

const entryCtxKey = progressKey(0)

type progressEntry struct {
	current Progress

	clbs []Callback

	parent *progressEntry

	children []*progressEntry
}

type ProgressTree struct {
	Progress

	Children []ProgressTree
}

func (e *progressEntry) progress() ProgressTree {
	p := ProgressTree{
		Progress: e.current,
		Children: make([]ProgressTree, 0, len(e.children)),
	}

	for _, ce := range e.children {
		prg := ce.progress()
		if prg.Progress == nil {
			continue
		}

		p.Children = append(p.Children, prg)
	}

	return p
}

func (e *progressEntry) addChild(child *progressEntry) {
	e.children = append(e.children, child)
}

func (e *progressEntry) update() {
	if e.parent != nil {
		e.parent.update()
	}

	for i := len(e.clbs) - 1; i >= 0; i-- {
		e.clbs[i](e.progress())
	}
}
