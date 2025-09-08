package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// DataItem represents an item flowing through the pipeline
type DataItem struct {
	ID    int
	Value string
	Data  interface{}
}

// Stage represents a pipeline stage
type Stage func(<-chan DataItem) <-chan DataItem

// Pipeline represents a processing pipeline
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates a new pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		stages: make([]Stage, 0),
	}
}

// AddStage adds a stage to the pipeline
func (p *Pipeline) AddStage(stage Stage) *Pipeline {
	p.stages = append(p.stages, stage)
	return p
}

// Execute runs the pipeline
func (p *Pipeline) Execute(input <-chan DataItem) <-chan DataItem {
	current := input
	for _, stage := range p.stages {
		current = stage(current)
	}
	return current
}

// Stage 1: Data validation and cleanup
func validationStage(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)
	go func() {
		defer close(output)
		for item := range input {
			fmt.Printf("Stage 1 - Validating item %d\n", item.ID)
			
			// Simulate validation work
			time.Sleep(100 * time.Millisecond)
			
			// Clean up the value (trim whitespace, convert to lowercase)
			item.Value = strings.TrimSpace(strings.ToLower(item.Value))
			
			// Only pass valid items (non-empty strings)
			if item.Value != "" {
				fmt.Printf("Stage 1 - Item %d validated: %s\n", item.ID, item.Value)
				output <- item
			} else {
				fmt.Printf("Stage 1 - Item %d rejected (empty value)\n", item.ID)
			}
		}
	}()
	return output
}

// Stage 2: Data transformation
func transformationStage(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)
	go func() {
		defer close(output)
		for item := range input {
			fmt.Printf("Stage 2 - Transforming item %d\n", item.ID)
			
			// Simulate transformation work
			time.Sleep(150 * time.Millisecond)
			
			// Try to convert string to number and calculate square
			if num, err := strconv.ParseFloat(item.Value, 64); err == nil {
				item.Data = math.Pow(num, 2)
				fmt.Printf("Stage 2 - Item %d transformed: %s -> %.2f\n", item.ID, item.Value, item.Data)
			} else {
				// If not a number, reverse the string
				reversed := reverseString(item.Value)
				item.Data = reversed
				fmt.Printf("Stage 2 - Item %d transformed: %s -> %s\n", item.ID, item.Value, reversed)
			}
			
			output <- item
		}
	}()
	return output
}

// Stage 3: Data enrichment
func enrichmentStage(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)
	go func() {
		defer close(output)
		for item := range input {
			fmt.Printf("Stage 3 - Enriching item %d\n", item.ID)
			
			// Simulate enrichment work
			time.Sleep(200 * time.Millisecond)
			
			// Create enriched data structure
			enrichedData := map[string]interface{}{
				"id":            item.ID,
				"original":      item.Value,
				"processed":     item.Data,
				"timestamp":     time.Now().Unix(),
				"processing_ms": 450, // Total processing time
			}
			
			item.Data = enrichedData
			fmt.Printf("Stage 3 - Item %d enriched with metadata\n", item.ID)
			output <- item
		}
	}()
	return output
}

// Stage 4: Final processing and formatting
func formattingStage(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)
	go func() {
		defer close(output)
		for item := range input {
			fmt.Printf("Stage 4 - Formatting item %d\n", item.ID)
			
			// Simulate formatting work
			time.Sleep(50 * time.Millisecond)
			
			// Format the final output
			if enrichedData, ok := item.Data.(map[string]interface{}); ok {
				formatted := fmt.Sprintf("Result[%d]: %s -> %v (processed at %d)",
					enrichedData["id"],
					enrichedData["original"],
					enrichedData["processed"],
					enrichedData["timestamp"])
				item.Data = formatted
			}
			
			fmt.Printf("Stage 4 - Item %d formatted\n", item.ID)
			output <- item
		}
	}()
	return output
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Data generator
func generateData() <-chan DataItem {
	output := make(chan DataItem)
	
	// Sample data
	testData := []string{
		"123.45",
		"hello world",
		"42",
		"",        // This should be filtered out
		"golang",
		"3.14159",
		"pipeline",
		"999",
	}
	
	go func() {
		defer close(output)
		for i, value := range testData {
			item := DataItem{
				ID:    i + 1,
				Value: value,
			}
			fmt.Printf("Generated item %d: %s\n", item.ID, item.Value)
			output <- item
			time.Sleep(50 * time.Millisecond) // Simulate data arrival rate
		}
	}()
	
	return output
}

func main() {
	fmt.Println("Starting Pipeline Processing Demo")
	fmt.Println(strings.Repeat("=", 50))
	
	// Create pipeline
	pipeline := NewPipeline().
		AddStage(validationStage).
		AddStage(transformationStage).
		AddStage(enrichmentStage).
		AddStage(formattingStage)
	
	// Generate input data
	input := generateData()
	
	// Execute pipeline
	output := pipeline.Execute(input)
	
	// Collect results
	fmt.Println("\nFinal Results:")
	fmt.Println(strings.Repeat("-", 30))
	
	var results []DataItem
	for result := range output {
		results = append(results, result)
		fmt.Printf("✓ %s\n", result.Data)
	}
	
	fmt.Printf("\nProcessing completed! Processed %d items successfully.\n", len(results))
}
