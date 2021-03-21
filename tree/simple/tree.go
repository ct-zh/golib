package simple

type node struct {
	val         int
	Left, Right *node
}

func NewNode(val int) *node {
	return &node{val, nil, nil}
}

func (n *node) SetVal(val int) {
	n.val = val
}

func (n *node) GetVal() int {
	return n.val
}

func (n *node) TraverseFn(f func(node2 *node)) {
	if n == nil {
		return
	}

	if n.Left != nil {
		n.Left.TraverseFn(f)
	}
	f(n)
	if n.Right != nil {
		n.Right.TraverseFn(f)
	}
}

func (n *node) TraverseWithChannel() chan *node {
	out := make(chan *node)
	go func() {
		n.TraverseFn(func(node2 *node) {
			out <- node2
		})
		close(out)
	}()
	return out
}
