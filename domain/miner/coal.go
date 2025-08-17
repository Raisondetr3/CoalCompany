package miner

type Coal struct {
	Count int
}

func NewCoal() *Coal {
	return &Coal{Count: 0}
}
