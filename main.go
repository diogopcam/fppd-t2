package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	// Configurações dos testes
	sizes := []int{1000, 10000, 50000}      // Tamanhos dos arrays
	granularities := []int{100, 500, 1000}  // Valores de G para MergeSort
	routinesList := []int{1, 2, 4}          // Valores de rotinas para InsertSort
	procsList := []int{1, 2, 4}             // Números de processadores

	// Inicialização
	rand.Seed(time.Now().UnixNano())

	// 1. Benchmark do MergeSort
	runMergeSortBenchmark(sizes, granularities, procsList, "resultados_merge.csv")

	// 2. Benchmark do InsertSort
	runInsertSortBenchmark(sizes, routinesList, procsList, "resultados_insert.csv")

	fmt.Println("Todos os benchmarks foram concluídos!")
}

// Funções para MergeSort
func runMergeSortBenchmark(sizes []int, granularities []int, procsList []int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo %s: %v\n", filename, err)
		return
	}
	defer file.Close()
	file.WriteString("Algoritmo,Tamanho,Granularidade,Processadores,Tempo\n")

	for _, size := range sizes {
		arr := generateRandomArray(size)
		for _, G := range granularities {
			// Tempo sequencial
			tSeq := benchmarkMergeSort(arr, G, 1, false)
			file.WriteString(fmt.Sprintf("mergesort,%d,%d,1,%.6f\n", size, G, tSeq.Seconds()))

			// Tempos paralelos
			for _, P := range procsList {
				arrCopy := make([]int, size)
				copy(arrCopy, arr)
				tPar := benchmarkMergeSort(arrCopy, G, P, true)
				file.WriteString(fmt.Sprintf("mergesort,%d,%d,%d,%.6f\n", size, G, P, tPar.Seconds()))
			}
		}
	}
	fmt.Printf("Benchmark MergeSort concluído. Resultados em %s\n", filename)
}

func benchmarkMergeSort(arr []int, G int, P int, useParallel bool) time.Duration {
	runtime.GOMAXPROCS(P)
	start := time.Now()

	if useParallel {
		_ = MergeSortParallel(arr, G)
	} else {
		_ = MergeSortSequential(arr)
	}

	return time.Since(start)
}

// Funções para InsertSort
func runInsertSortBenchmark(sizes []int, routinesList []int, procsList []int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo %s: %v\n", filename, err)
		return
	}
	defer file.Close()
	file.WriteString("Algoritmo,Tamanho,NumRotinas,Processadores,Tempo\n")

	for _, size := range sizes {
		arr := generateRandomArray(size)
		// Tempo sequencial
		tSeq := benchmarkInsertSort(arr, 1, 1, false)
		file.WriteString(fmt.Sprintf("insertsort,%d,1,1,%.6f\n", size, tSeq.Seconds()))

		// Tempos paralelos
		for _, routines := range routinesList {
			for _, P := range procsList {
				arrCopy := make([]int, size)
				copy(arrCopy, arr)
				tPar := benchmarkInsertSort(arrCopy, routines, P, true)
				file.WriteString(fmt.Sprintf("insertsort,%d,%d,%d,%.6f\n", size, routines, P, tPar.Seconds()))
			}
		}
	}
	fmt.Printf("Benchmark InsertSort concluído. Resultados em %s\n", filename)
}

func benchmarkInsertSort(arr []int, routines int, P int, useParallel bool) time.Duration {
	runtime.GOMAXPROCS(P)
	start := time.Now()

	if useParallel {
		InsertionSortParallel(arr, routines)
	} else {
		InsertionSortSequential(arr)
	}

	return time.Since(start)
}

// Função auxiliar comum
func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size * 10)
	}
	return arr
}