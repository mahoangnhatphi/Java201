package behavioral

import "fmt"

// SortStrategy defines the interface for sorting strategies
type SortStrategy interface {
	Sort(data []int) []int
}

// BubbleSort implements SortStrategy
type BubbleSort struct{}

func NewBubbleSort() *BubbleSort {
	return &BubbleSort{}
}

func (b *BubbleSort) Sort(data []int) []int {
	arr := make([]int, len(data))
	copy(arr, data)
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// QuickSort implements SortStrategy
type QuickSort struct{}

func NewQuickSort() *QuickSort {
	return &QuickSort{}
}

func (q *QuickSort) Sort(data []int) []int {
	arr := make([]int, len(data))
	copy(arr, data)
	quickSort(arr, 0, len(arr)-1)
	return arr
}

func quickSort(arr []int, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// MergeSort implements SortStrategy
type MergeSort struct{}

func NewMergeSort() *MergeSort {
	return &MergeSort{}
}

func (m *MergeSort) Sort(data []int) []int {
	arr := make([]int, len(data))
	copy(arr, data)
	return mergeSort(arr)
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// Sorter is the context that uses a strategy
type Sorter struct {
	strategy SortStrategy
}

func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

func (s *Sorter) Sort(data []int) []int {
	return s.strategy.Sort(data)
}

// StrategyExampleUsage demonstrates the Strategy pattern
func StrategyExampleUsage() {
	data := []int{64, 34, 25, 12, 22, 11, 90}

	sorter := NewSorter(NewBubbleSort())
	fmt.Printf("Bubble Sort: %v\n", sorter.Sort(data))

	sorter.SetStrategy(NewQuickSort())
	fmt.Printf("Quick Sort: %v\n", sorter.Sort(data))

	sorter.SetStrategy(NewMergeSort())
	fmt.Printf("Merge Sort: %v\n", sorter.Sort(data))
}
