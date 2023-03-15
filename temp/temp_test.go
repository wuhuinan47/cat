package temp

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	matrix := [][]int{
		{-1, -1, -1, 1, 0, 0, 0},
		{-1, -1, -1, 0, 0, 0, 5},
		{2, 0, -1, 0, 3, 0, 0},
		{0, -1, -1, -1, 0, 1, 5},
		{0, -1, -1, -1, 0, 0, 0},
		{1, -1, -1, -1, 2, 0, 3},
		{-1, -1, -1, -1, -1, -1, -1},
		{4, 0, 0, 0, 0, 0, 0},
	}

	list, index := findLastNegativeList(matrix)

	fmt.Printf("Last negative list is %v, and its index is %d, ac is %v, haha is %v\n", list, index, matrix[index], formatNextRows(list, matrix[index+1]))
}

// 格式化下一行的列表,1表示可挖 0表示不可挖
func formatNextRows(list, list2 []int) []int {
	len := len(list)
	xx := make([]int, len)

	for i := 0; i < len; i++ {
		if list[i] == -1 {
			xx[i] = 1
		} else {
			xx[i] = 0
		}
	}

	return xx
}

// 找到最后一个包含-1的列表,并返回该列表和其索引
func findLastNegativeList(matrix [][]int) ([]int, int) {
	index := -1
	for i := len(matrix) - 1; i >= 0; i-- {
		if containsNegative(matrix[i]) {
			index = i
			break
		}
	}
	return matrix[index], index
}

// 若列表中包含-1，则返回true
func containsNegative(list []int) bool {
	for _, value := range list {
		if value == -1 {
			return true
		}
	}
	return false
}
