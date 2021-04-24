package gsort

// 时间复杂读O(n^2)
// 在[0, l-i-1)的区间内元素两两对比，小的放前面大的放后面
// 两两比较，i 与 i + 1比较， 所以i和j的结尾都需要-1
func bubble(a []int) {
	for i := 0; i < len(a)-1; i++ {
		for j := 0; j < len(a)-i-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
			}
		}
	}
}
