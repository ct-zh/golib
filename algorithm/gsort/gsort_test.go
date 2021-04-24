package gsort

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var sortMap = map[SortMethod]string{
	BubbleSort:    "冒泡排序",
	SelectSort:    "选择排序",
	InsertionSort: "插入排序",
	QuickSort:     "快速排序",
	MergeSort:     "归并排序",
}

func TestSetSort(t *testing.T) {
	type args struct {
		a      []int
		method SortMethod
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{
				a: []int{5, 4, 3, 2, 1},
			},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for method, name := range sortMap {
				SetSort(tt.args.a, method)
				if !CheckSort(tt.args.a) {
					t.Fatalf("%s有问题：结果：%+v", name, tt.args.a)
				}
				if !reflect.DeepEqual(tt.want, tt.args.a) {
					t.Fatalf("%s有问题：结果：%+v", name, tt.args.a)
				}
				t.Logf("%s ok！", name)
			}
		})
	}
}

// 生成元素数量为n,元素在[rangeL,rangeR]区间的 整数集合
func GenerateRandomArray(n int, rangeL int, rangeR int) []int {
	if rangeL > rangeR {
		return nil
	}

	arr := make([]int, n) // rand

	rand2 := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < n; i++ {
		// 指定范围随机数生成的标准写法
		arr[i] = rand2.Int()%(rangeR-rangeL+1) + rangeL
	}

	return arr
}

// 生成近似顺序的整数数组
// n: 数组大小; swapTimes: 数组打乱次数;
func GenerateNearlyArray(n int, swapTimes int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}
	rand2 := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < swapTimes; i++ {
		// 取两个 [0,n) 区间的数
		posX := rand2.Int() % n
		posy := rand2.Int() % n
		// swap 随机交换两个数据
		arr[posX], arr[posy] = arr[posy], arr[posX]
	}

	return arr
}

// 检查数组是否按照升序排序的
func CheckSort(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}
