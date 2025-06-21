package main

import (
	"runtime"
	"sync"
	"time"
)

func InsertionSortSequential(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func InsertionSortParallel(arr []int, numRoutines int) {
	if numRoutines <= 1 || len(arr) < numRoutines*2 {
		InsertionSortSequential(arr)
		return
	}

	size := len(arr) / numRoutines
	var wg sync.WaitGroup

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		start := i * size
		end := start + size
		if i == numRoutines-1 {
			end = len(arr)
		}

		go func(s, e int) {
			defer wg.Done()
			InsertionSortSequential(arr[s:e])
		}(start, end)
	}
	wg.Wait()

	for i := 1; i < numRoutines; i++ {
		start := 0
		end := i * size
		if i == numRoutines-1 {
			end = len(arr)
		}
		mergeInsertion(arr, start, end)
	}
}

func mergeInsertion(arr []int, start, end int) {
	for i := start + 1; i < end; i++ {
		key := arr[i]
		j := i - 1

		for j >= start && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func BenchmarkInsertionSort(arr []int, numRoutines int, parallel bool) time.Duration {
	start := time.Now()

	if parallel {
		runtime.GOMAXPROCS(numRoutines)
		InsertionSortParallel(arr, numRoutines)
	} else {
		InsertionSortSequential(arr)
	}

	return time.Since(start)
}