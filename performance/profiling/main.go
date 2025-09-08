package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// PerformanceTracker tracks execution time and memory usage
type PerformanceTracker struct {
	Name      string
	StartTime time.Time
	StartMem  runtime.MemStats
}

// NewPerformanceTracker creates a new performance tracker
func NewPerformanceTracker(name string) *PerformanceTracker {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return &PerformanceTracker{
		Name:      name,
		StartTime: time.Now(),
		StartMem:  memStats,
	}
}

// End stops tracking and prints results
func (pt *PerformanceTracker) End() {
	duration := time.Since(pt.StartTime)
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	memUsed := memStats.TotalAlloc - pt.StartMem.TotalAlloc
	
	fmt.Printf("📊 [%s] Duration: %v, Memory: %d bytes\n", pt.Name, duration, memUsed)
}

// Inefficient string concatenation
func inefficientStringConcat(n int) string {
	tracker := NewPerformanceTracker("Inefficient String Concat")
	defer tracker.End()
	
	result := ""
	for i := 0; i < n; i++ {
		result += fmt.Sprintf("item_%d ", i)
	}
	return result
}

// Efficient string concatenation using StringBuilder
func efficientStringConcat(n int) string {
	tracker := NewPerformanceTracker("Efficient String Concat")
	defer tracker.End()
	
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(fmt.Sprintf("item_%d ", i))
	}
	return builder.String()
}

// Inefficient slice operations
func inefficientSliceGrowth(n int) []int {
	tracker := NewPerformanceTracker("Inefficient Slice Growth")
	defer tracker.End()
	
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

// Efficient slice operations with pre-allocation
func efficientSliceGrowth(n int) []int {
	tracker := NewPerformanceTracker("Efficient Slice Growth")
	defer tracker.End()
	
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return result
}

// CPU-intensive task without optimization
func cpuIntensiveTask(data []int) int {
	tracker := NewPerformanceTracker("CPU Intensive (Serial)")
	defer tracker.End()
	
	sum := 0
	for _, v := range data {
		// Simulate expensive calculation
		for i := 0; i < 1000; i++ {
			sum += v * i
		}
	}
	return sum
}

// CPU-intensive task with goroutine optimization
func cpuIntensiveTaskParallel(data []int) int {
	tracker := NewPerformanceTracker("CPU Intensive (Parallel)")
	defer tracker.End()
	
	numWorkers := runtime.NumCPU()
	chunkSize := len(data) / numWorkers
	
	results := make(chan int, numWorkers)
	
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(data)
		}
		
		go func(chunk []int) {
			sum := 0
			for _, v := range chunk {
				for j := 0; j < 1000; j++ {
					sum += v * j
				}
			}
			results <- sum
		}(data[start:end])
	}
	
	totalSum := 0
	for i := 0; i < numWorkers; i++ {
		totalSum += <-results
	}
	
	return totalSum
}

// Memory-intensive operations
type DataPoint struct {
	ID        int
	Value     float64
	Timestamp time.Time
	Metadata  map[string]string
}

// Inefficient memory usage with lots of allocations
func inefficientMemoryUsage(n int) []*DataPoint {
	tracker := NewPerformanceTracker("Inefficient Memory Usage")
	defer tracker.End()
	
	var data []*DataPoint
	for i := 0; i < n; i++ {
		point := &DataPoint{
			ID:        i,
			Value:     rand.Float64(),
			Timestamp: time.Now(),
			Metadata:  make(map[string]string),
		}
		point.Metadata["source"] = "sensor"
		point.Metadata["location"] = fmt.Sprintf("loc_%d", i%10)
		point.Metadata["type"] = "temperature"
		
		data = append(data, point)
	}
	return data
}

// More efficient memory usage with object pooling
type DataPointPool struct {
	pool sync.Pool
}

func NewDataPointPool() *DataPointPool {
	return &DataPointPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &DataPoint{
					Metadata: make(map[string]string),
				}
			},
		},
	}
}

func (dp *DataPointPool) Get() *DataPoint {
	return dp.pool.Get().(*DataPoint)
}

func (dp *DataPointPool) Put(point *DataPoint) {
	// Clear the data
	point.ID = 0
	point.Value = 0
	point.Timestamp = time.Time{}
	for k := range point.Metadata {
		delete(point.Metadata, k)
	}
	dp.pool.Put(point)
}

func efficientMemoryUsage(n int) []*DataPoint {
	tracker := NewPerformanceTracker("Efficient Memory Usage (Pool)")
	defer tracker.End()
	
	pool := NewDataPointPool()
	data := make([]*DataPoint, 0, n)
	
	for i := 0; i < n; i++ {
		point := pool.Get()
		point.ID = i
		point.Value = rand.Float64()
		point.Timestamp = time.Now()
		point.Metadata["source"] = "sensor"
		point.Metadata["location"] = fmt.Sprintf("loc_%d", i%10)
		point.Metadata["type"] = "temperature"
		
		data = append(data, point)
	}
	
	// In a real scenario, you would return objects to pool when done
	return data
}

// Sorting performance comparison
func compareSort(data []int) {
	// Copy data for fair comparison
	data1 := make([]int, len(data))
	data2 := make([]int, len(data))
	data3 := make([]int, len(data))
	copy(data1, data)
	copy(data2, data)
	copy(data3, data)
	
	// Standard library sort
	tracker1 := NewPerformanceTracker("Standard Sort")
	sort.Ints(data1)
	tracker1.End()
	
	// Bubble sort (inefficient)
	tracker2 := NewPerformanceTracker("Bubble Sort")
	bubbleSort(data2)
	tracker2.End()
	
	// Quick sort implementation
	tracker3 := NewPerformanceTracker("Quick Sort")
	quickSort(data3, 0, len(data3)-1)
	tracker3.End()
}

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
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
	
	for j := low; j <= high-1; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// Benchmark function
func runBenchmark(name string, fn func(), iterations int) {
	fmt.Printf("\n🏃 Running benchmark: %s\n", name)
	fmt.Println(strings.Repeat("=", 40))
	
	start := time.Now()
	var totalMem runtime.MemStats
	runtime.ReadMemStats(&totalMem)
	startAlloc := totalMem.TotalAlloc
	
	for i := 0; i < iterations; i++ {
		fn()
	}
	
	duration := time.Since(start)
	runtime.ReadMemStats(&totalMem)
	memUsed := totalMem.TotalAlloc - startAlloc
	
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average time: %v\n", duration/time.Duration(iterations))
	fmt.Printf("Total memory: %d bytes\n", memUsed)
	fmt.Printf("Average memory: %d bytes\n", memUsed/uint64(iterations))
}

func demonstrateStringPerformance() {
	fmt.Println("\n📝 String Concatenation Performance")
	fmt.Println(strings.Repeat("-", 40))
	
	n := 1000
	inefficientStringConcat(n)
	efficientStringConcat(n)
}

func demonstrateSlicePerformance() {
	fmt.Println("\n🔢 Slice Operations Performance")
	fmt.Println(strings.Repeat("-", 40))
	
	n := 100000
	inefficientSliceGrowth(n)
	efficientSliceGrowth(n)
}

func demonstrateCPUPerformance() {
	fmt.Println("\n💻 CPU Intensive Tasks Performance")
	fmt.Println(strings.Repeat("-", 40))
	
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(100)
	}
	
	cpuIntensiveTask(data)
	cpuIntensiveTaskParallel(data)
}

func demonstrateMemoryPerformance() {
	fmt.Println("\n🧠 Memory Usage Performance")
	fmt.Println(strings.Repeat("-", 40))
	
	n := 10000
	inefficientMemoryUsage(n)
	efficientMemoryUsage(n)
}

func demonstrateSortingPerformance() {
	fmt.Println("\n📊 Sorting Algorithms Performance")
	fmt.Println(strings.Repeat("-", 40))
	
	// Generate random data
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(1000)
	}
	
	compareSort(data)
}

func printSystemInfo() {
	fmt.Println("🖥️  System Information")
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Go Version: %s\n", runtime.Version())
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Current Alloc: %d KB\n", m.Alloc/1024)
	fmt.Printf("Total Alloc: %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys: %d KB\n", m.Sys/1024)
}

func main() {
	fmt.Println("Performance Analysis and Optimization Demo")
	fmt.Println("=========================================")
	
	printSystemInfo()
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Run different performance demonstrations
	demonstrateStringPerformance()
	demonstrateSlicePerformance()
	demonstrateCPUPerformance()
	demonstrateMemoryPerformance()
	demonstrateSortingPerformance()
	
	fmt.Println("\n✅ Performance analysis completed!")
	fmt.Println("Key takeaways:")
	fmt.Println("1. Pre-allocate slices when size is known")
	fmt.Println("2. Use strings.Builder for string concatenation")
	fmt.Println("3. Leverage goroutines for CPU-intensive parallel tasks")
	fmt.Println("4. Consider object pooling for frequent allocations")
	fmt.Println("5. Choose appropriate algorithms for your use case")
}
