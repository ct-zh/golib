package gsort

// 自选排序，为了与sort包区分，加了g前缀

type SortMethod int

const (
	SelectSort SortMethod = iota + 1
	BubbleSort
	InsertionSort
	MergeSort
	QuickSort
)

func SetSort(a []int, method SortMethod) {
	switch method {
	case SelectSort:
		selectSort(a)
	case BubbleSort:
		bubble(a)
	case InsertionSort:
		insertion(a)
	case MergeSort:
		mergeSort(a)
	case QuickSort:
		quick(a)
	}
}
