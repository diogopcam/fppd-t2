package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	// Configuração inicial
	maxProcs := runtime.NumCPU() // Detecta automaticamente (10 no seu caso)
	fmt.Printf("Utilizando até %d núcleos físicos\n", maxProcs)

	// Configurações dos testes otimizadas para 10 núcleos
	sizes := []int{10000, 50000, 100000, 500000}      // Tamanhos maiores para aproveitar paralelismo
	granularities := []int{500, 1000, 5000}           // Valores de G para MergeSort
	routinesList := []int{1, 5, 10, 15, 20}           // Valores de rotinas para InsertSort
	procsList := []int{1, 2, 4, 6, 8, 10}             // Progressão até todos os núcleos

	// Aquecimento do sistema
	warmUp(maxProcs)

	// Execução dos benchmarks
	runMergeSortBenchmark(sizes, granularities, procsList, "resultados_merge.csv")
	runInsertSortBenchmark(sizes, routinesList, procsList, "resultados_insert.csv")

	fmt.Println("Benchmarks concluídos com sucesso!")
}

func warmUp(maxProcs int) {
	fmt.Println("Aquecendo o sistema...")
	warmArr := make([]int, 100000)
	rand.Seed(time.Now().UnixNano())
	
	// Aquecimento do MergeSort
	runtime.GOMAXPROCS(maxProcs)
	_ = MergeSortParallel(warmArr, 1000)
	
	// Aquecimento do InsertSort
	InsertionSortParallel(warmArr, maxProcs*2)
}

// Funções para MergeSort (otimizadas)
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
			// Teste sequencial
			tSeq := benchmarkMergeSort(arr, G, 1, false)
			recordResult(file, "mergesort", size, G, 1, tSeq)

			// Testes paralelos
			for _, P := range procsList {
				arrCopy := make([]int, size)
				copy(arrCopy, arr)
				tPar := benchmarkMergeSort(arrCopy, G, P, true)
				recordResult(file, "mergesort", size, G, P, tPar)
				
				// Pausa para evitar throttling
				if P >= 6 {
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}
	fmt.Printf("Benchmark MergeSort salvo em %s\n", filename)
}

func benchmarkMergeSort(arr []int, G int, P int, parallel bool) time.Duration {
	runtime.GOMAXPROCS(P)
	start := time.Now()

	// Ajuste dinâmico da granularidade
	if parallel {
		optimalG := len(arr) / (P * 4)
		if optimalG < G {
			G = optimalG
		}
		_ = MergeSortParallel(arr, G)
	} else {
		_ = MergeSortSequential(arr)
	}

	return time.Since(start)
}

// Funções para InsertSort (otimizadas)
func runInsertSortBenchmark(sizes []int, routinesList []int, procsList []int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo %s: %v\n", filename, err)
		return
	}
	defer file.Close()
	file.WriteString("Algoritmo,Tamanho,Rotinas,Processadores,Tempo\n")

	for _, size := range sizes {
		arr := generateRandomArray(size)
		
		// Teste sequencial
		tSeq := benchmarkInsertSort(arr, 1, 1, false)
		recordResult(file, "insertsort", size, 1, 1, tSeq)

		// Testes paralelos
		for _, routines := range routinesList {
			for _, P := range procsList {
				arrCopy := make([]int, size)
				copy(arrCopy, arr)
				tPar := benchmarkInsertSort(arrCopy, routines, P, true)
				recordResult(file, "insertsort", size, routines, P, tPar)
				
				if P >= 6 {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}
	fmt.Printf("Benchmark InsertSort salvo em %s\n", filename)
}

func benchmarkInsertSort(arr []int, routines int, P int, parallel bool) time.Duration {
	runtime.GOMAXPROCS(P)
	start := time.Now()

	if parallel {
		// Otimização: número de rotinas baseado nos núcleos
		optimalRoutines := P * 2
		if routines > optimalRoutines {
			optimalRoutines = routines
		}
		InsertionSortParallel(arr, optimalRoutines)
	} else {
		InsertionSortSequential(arr)
	}

	return time.Since(start)
}

// Funções auxiliares
func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size * 10)
	}
	return arr
}

func recordResult(file *os.File, algorithm string, size, param, procs int, duration time.Duration) {
	var line string
	if algorithm == "mergesort" {
		line = fmt.Sprintf("%s,%d,%d,%d,%.6f\n", algorithm, size, param, procs, duration.Seconds())
	} else {
		line = fmt.Sprintf("%s,%d,%d,%d,%.6f\n", algorithm, size, param, procs, duration.Seconds())
	}
	file.WriteString(line)
	fmt.Print(line) // Exibe progresso em tempo real
}