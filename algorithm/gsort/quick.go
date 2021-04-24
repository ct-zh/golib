package gsort

// 快速排序
// 对[1, len)的数据做判断，保证小于a[0]的在左边，大于a[0]的在右边，
// 此时a[0]移到中间成为分割该区间的partition
// 再对子数组partition，直到最小粒度，完成排序

func quick(a []int) {
	quick2(0, len(a)-1, a) // [0, len(a) - 1]
}

func quick2(l, r int, a []int) {
	if l >= r {
		return
	}

	lt, gt := partition(l, r, a) // [l, r]
	quick2(l, lt-1, a)           // [l, lt-1]
	quick2(gt, r, a)             // [gt, r]
}

// 在做最后一步a[0]与a[lt]交换之前，都保证
// [l+1, lt] < [lt+1, gt-1] < [gt, r]
// 交换之后变成
// [l, lt-1] < [lt, gt - 1] < [gt, r]
func partition(l, r int, a []int) (lt, gt int) {
	g := a[l]

	//  [l+1, r]
	lt, gt = l, r+1

	i := l + 1
	for i < gt {
		if a[i] == g {
			i++
		} else if a[i] < g {
			a[lt+1], a[i] = a[i], a[lt+1]
			lt++
			i++
		} else {
			a[gt-1], a[i] = a[i], a[gt-1]
			gt--
		}
	}
	a[l], a[lt] = a[lt], a[l]

	return
}
