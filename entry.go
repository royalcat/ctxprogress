package ctxprogress

type progressKey uint

const entryCtxKey = progressKey(0)

type progressEntry struct {
	current Progress

	clbsTrs []CallbackTree
	clbs    []Callback

	parent   *progressEntry
	children []*progressEntry
}

type ProgressTree struct {
	Entry    Progress
	Children []ProgressTree
}

// Progress implements Progress.
func (p ProgressTree) Progress() (current int, total int) {
	return p.Entry.Progress()
}

func (e *progressEntry) progress() ProgressTree {
	p := ProgressTree{
		Entry:    e.current,
		Children: make([]ProgressTree, 0, len(e.children)),
	}

	for _, ce := range e.children {
		prg := ce.progress()
		if prg.Entry == nil {
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

	if len(e.clbsTrs) > 0 || len(e.clbs) > 0 {
		pt := e.progress()
		for i := len(e.clbsTrs) - 1; i >= 0; i-- {
			e.clbsTrs[i](pt)
		}

		for i := len(e.clbs) - 1; i >= 0; i-- {
			e.clbs[i](pt.Entry)
		}
	}

}
