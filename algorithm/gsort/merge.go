package gsort

// 归并排序
// 将数组不停对半切分，直到不能切分为止；
// 然后再两两对比合并 (递归)

// 自顶向下的归并排序
func mergeSort(a []int) {
	_mergeSort(0, len(a)-1, a)
}

func _mergeSort(l, r int, a []int) {
	if l >= r {
		return
	}

	mid := (r-l)/2 + l // 做减法，防止溢出

	_mergeSort(l, mid, a)   // [l, mid]
	_mergeSort(mid+1, r, a) // (mid, r]

	merge(l, mid, r, a)
}

// 归并两个区间 [l, mid] (mid, r]
// 方法应该是双指针两两比较？
func merge(l, mid, r int, a []int) {
	tmp := make([]int, r-l+1) // [l, r] 闭区间

	i1 := l
	i2 := mid + 1
	i := 0
	for {
		if i >= cap(tmp) {
			break
		}

		if i1 <= mid && i2 <= r {
			if a[i1] < a[i2] {
				tmp[i] = a[i1]
				i1++
			} else {
				tmp[i] = a[i2]
				i2++
			}
		} else if i1 <= mid {
			tmp[i] = a[i1]
			i1++
		} else {
			tmp[i] = a[i2]
			i2++
		}
		i++
	}

	for k, v := range tmp {
		a[l+k] = v
	}
}

// 自底向上的归并排序
// 也就是第一步就将数组分成最小块，然后依次进行归并操作
// 性能对比自顶向下的归并排序速度可能会稍慢一些 => 为什么?
// 但是可以看到代码里面没有用到数组的key值，也就是说这个算法可以用在链表上
func mergeSort1(arr []int) {
	for size := 1; size <= len(arr); size += size { // 模块大小 1 2 4 8 ... 直到和数组长度相等，即完成所有归并
		for i := 0; i+size < len(arr); i += size + size { // 每次都是两个模块做比较，所以自增 2size
			// 注意处理越界问题：
			// 对于一个模块来说,内部已经是排好序了,所以每次都是两个模块进行比较,我们需要保证这两个模块的边界
			// 第一个模块[i, i+size-1],需要保证: i+size 小于等于Count (在for循环的条件里保证了)
			// 第二个模块[i+size, i+size+size-1] ，需要保证 endKey 不能大于count-1
			// (因为第二个模块不需要是完整的长度为size的模块，允许长度不足size，所以不放在for循环的条件里限制死)
			merge(i, i+size-1,
				minInt(i+size+size-1, len(arr)-1), // 保证 endKey 不能大于count-1
				arr)
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
