package ctxprogress

type progressKey uint

const entryCtxKey = progressKey(0)

type progressEntry[P any] struct {
	total      int
	current    int
	properties P

	parent *progressEntry[P]
	clbs   []Callback[P]

	children []*progressEntry[P]
}

func (e *progressEntry[D]) progress() Progress[D] {
	p := Progress[D]{
		Properties: e.properties,
		Current:    e.current,
		Total:      e.total,
		Children:   make([]Progress[D], 0, len(e.children)),
	}

	for _, ce := range e.children {
		p.Children = append(p.Children, ce.progress())
	}

	return p
}

func (e *progressEntry[D]) addChild(child *progressEntry[D]) {
	e.children = append(e.children, child)
}

func (e *progressEntry[D]) update() {
	if e.parent != nil {
		e.parent.update()
	}

	for i := len(e.clbs) - 1; i >= 0; i-- {
		e.clbs[i](e.progress())
	}
}
