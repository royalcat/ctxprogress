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

func sumProgressTree(c []ProgressTree) (current int, total int) {
	for _, p := range c {
		c, t := p.Progress()
		current += c
		total += t
	}

	return current, total

}

// Progress implements Progress.
func (p ProgressTree) Progress() (current int, total int) {
	if p.Entry == nil {
		return sumProgressTree(p.Children)
	}

	if len(p.Children) > 0 {
		ec, et := p.Entry.Progress()
		if ec > len(p.Children)-1 {
			return ec, et
		}

		lc, lt := p.Children[ec].Progress()
		current = ec * lt
		if ec != et {
			current += lc
		}

		total = et * lt

		return current, total
	}

	return p.Entry.Progress()
}

func (e *progressEntry) progress() ProgressTree {
	p := ProgressTree{
		Entry:    e.current,
		Children: make([]ProgressTree, 0, len(e.children)),
	}

	for _, ce := range e.children {
		p.Children = append(p.Children, ce.progress())
	}

	return p
}

func (e *progressEntry) addChild(child *progressEntry) {
	e.children = append(e.children, child)
}

func (e *progressEntry) update() {
	if len(e.clbsTrs) > 0 || len(e.clbs) > 0 {
		pt := e.progress()
		for i := len(e.clbsTrs) - 1; i >= 0; i-- {
			e.clbsTrs[i](pt)
		}

		for i := len(e.clbs) - 1; i >= 0; i-- {
			e.clbs[i](pt)
		}
	}

	if e.parent != nil {
		e.parent.update()
	}

}
