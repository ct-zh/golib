package gsort

// 选择排序 时间复杂度 O(n^2)
// 每次在子list [i+1, len) 中找出小于i的最小值min，与i交换
func selectSort(a []int) {
	for i := 0; i < len(a); i++ {
		min := i
		for j := i + 1; j < len(a); j++ {
			if a[min] > a[j] {
				min = j
			}
		}
		if min != i {
			a[i], a[min] = a[min], a[i]
		}
	}
}
