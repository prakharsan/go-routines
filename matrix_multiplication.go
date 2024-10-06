package main

import (
	"fmt"
	"sync"
)

func multiplyElement(A, B [][]int, result [][]int, i, j, N int, wg *sync.WaitGroup) {
	defer wg.Done()

	sum := 0
	for k := 0; k < N; k++ {
		sum += A[i][k] * B[k][j]
	}
	result[i][j] = sum
}

func main() {
	N := 3

	A := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	B := [][]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}

	result := make([][]int, N)
	for i := range result {
		result[i] = make([]int, N)
	}

	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			wg.Add(1)
			go multiplyElement(A, B, result, i, j, N, &wg)
		}
	}

	wg.Wait()

	fmt.Println("Result Matrix C:")
	for _, row := range result {
		fmt.Println(row)
	}
}
