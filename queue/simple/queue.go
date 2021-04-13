package simple

type queue interface {
	LPush(v interface{})
	LPop() interface{}
	RPush(v interface{})
	RPop() interface{}
	IsEmpty() bool
	Len() int
}

type simpleQueue []interface{}

func New() *simpleQueue {
	return &simpleQueue{}
}

func (q *simpleQueue) RPush(v interface{}) {
	*q = append(*q, v)
}

func (q *simpleQueue) RPop() interface{} {
	tail := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return tail
}

func (q *simpleQueue) LPush(v interface{}) {
	head := []interface{}{v}
	*q = append(head, *q...)
}

func (q *simpleQueue) LPop() interface{} {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *simpleQueue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *simpleQueue) Len() int {
	return len(*q)
}
