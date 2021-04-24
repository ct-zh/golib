# 排序

## 大纲

自写排序，通过g前缀与sort包区分；实现了比较常见的一些排序算法，利用工厂方法`SetSort`进行排序方法自由选择；

1. 选择排序
2. 插入排序与插入排序的优化(少进行一次数据交换)
3. 自顶向下的归并排序，自底向上的归并排序
4. 冒泡排序
5. 希尔排序
6. 快排,随机快排,双路,三路快排
7. 排序的性能比较
8. 归并排序的另外一个优化，在merge外申请aux空间

> [liuyubobo的玩转数据结构](https://github.com/liuyubobobo/Play-with-Algorithms)

## 时间复杂度比较
1. 选择排序，复杂度O(n^2)；
2. 插入排序，复杂度最差的时候是O(n^2)，对近似有序的数组则是无限接近与O(n)，因为判断过程中存在提前break的情况，所以性能优于选择排序；
3. 归并排序，复杂度O(nlogn)，在数据量小的时候使用插入排序；
4. 快速排序，复杂度O(nlogn)，对近似有序的数组则退化为O(n^2)，在数据量小的时候仍然建议使用插入排序；

```go
// 不同排序算法的性能比较, 平均时间复杂度
// 冒泡排序(bubble sort)    O(n^2)
// 选择排序(selection sort) O(n^2)
// 插入排序(insertion sort) O(n^2)
// 希尔排序(shell's sort)   O(n^1.5) 时间复杂度下界为 O(n*log2n)
// 快速排序(quick sort)     O(n*logN)
// 归并排序(merge sort)	  O(n*logN)
// 堆排序(heap sort)		  O(n*logN)
// 基数排序(radix sort)	  O(n*log(r)m) r为基数，m为堆数
```

