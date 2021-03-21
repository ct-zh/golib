package simple

type queue []interface{}

func New() *queue {
	return &queue{}
}

func (q *queue) Push(v interface{}) {
	*q = append(*q, v)
}

func (q *queue) Pop() interface{} {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *queue) IsEmpty() bool {
	return len(*q) == 0
}
