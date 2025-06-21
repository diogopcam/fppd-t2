package main

import "sync"

func MergeSortSequential(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    mid := len(arr) / 2
    left := MergeSortSequential(arr[:mid])
    right := MergeSortSequential(arr[mid:])
    return merge(left, right)
}

func MergeSortParallel(arr []int, G int) []int {
    if len(arr) <= G {
        return MergeSortSequential(arr)
    }
    
    mid := len(arr) / 2
    var left, right []int
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        left = MergeSortParallel(arr[:mid], G)
    }()
    right = MergeSortParallel(arr[mid:], G)
    wg.Wait()
    
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] < right[j] {
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