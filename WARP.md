# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Repository Overview

This is a Go demonstration repository containing independent examples and performance studies. The codebase is structured as a single Go module (`github.com/kenneth-wang/go-demo`) with each directory representing a separate demo or experiment.

## Architecture

### Module Structure
- **Package-per-directory**: Each subdirectory contains a focused Go package demonstrating specific concepts
- **Independent examples**: Each demo can be run independently without dependencies on other demos
- **Mixed executable types**: Contains both `main` packages (runnable) and library packages (testable)

### Key Components
- **ctxmemgrowth/**: Memory growth analysis for Go contexts using `context.WithValue()`
- **delaycall/**: Complex concurrency pattern implementing user-specific request delaying with goroutine pools
- **interface/**: Interface testing patterns with mocking using gomonkey
- **kafka/**: Kafka integration with SASL/SCRAM authentication for AWS MSK
- **concurrency/workerpool/**: Worker pool pattern with job queue and result collection
- **concurrency/pipeline/**: Pipeline pattern for data processing with multiple stages
- **datastructures/lru/**: LRU cache implementation with doubly linked list and hash map
- **patterns/observer/**: Observer pattern implementation with event dispatching system
- **performance/profiling/**: Performance analysis and optimization techniques demonstration
- **networking/tcp-echo/**: TCP echo server and client with multi-client support

## Common Development Commands

### Building and Running
```bash
# Run a specific demo
go run ./ctxmemgrowth/
go run ./kafka/

# Run delaycall (note: package name differs from directory)
go run ./delaycall/

# Run new demo cases
go run ./concurrency/workerpool/
go run ./concurrency/pipeline/
go run ./datastructures/lru/
go run ./patterns/observer/
go run ./performance/profiling/
go run ./networking/tcp-echo/

# Run TCP server in interactive mode
go run ./networking/tcp-echo/ interactive

# Build all packages to check compilation
go build ./...
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./delaycall/
go test ./interface/

# Run tests with verbose output
go test -v ./...

# Run specific test function
go test -run TestGetValueByKey_Normal ./interface/
```

### Dependencies
```bash
# Download and verify dependencies
go mod download
go mod verify

# Clean up unused dependencies
go mod tidy

# View dependency graph
go mod graph
```

## Development Notes

### Testing Patterns
- The codebase uses **gomonkey** for method mocking in tests (see `interface/main_test.go`)
- Tests include both normal flow and panic/error condition testing
- Long-running tests exist (delaycall test runs for extended periods)

### Concurrency Patterns
- **delaycall** package demonstrates advanced patterns:
  - Dynamic goroutine spawning per user
  - Channel-based request routing
  - Timer-based cleanup with graceful shutdown
  - Mutex-protected shared state management
- **concurrency/workerpool** shows job distribution patterns:
  - Worker pool with bounded goroutines
  - Job queue with result collection
  - Graceful shutdown coordination
- **concurrency/pipeline** demonstrates data flow patterns:
  - Multi-stage processing pipeline
  - Channel-based stage communication
  - Concurrent stage processing

### External Dependencies
Key dependencies include:
- **IBM/sarama**: Kafka client library
- **gomonkey/v2**: Runtime method patching for tests  
- **gin-gonic/gin**: Web framework (referenced in go.mod)
- **spf13/cobra**: CLI framework
- **google/uuid**: UUID generation

### Data Structures
- **datastructures/lru** implements efficient LRU cache:
  - Doubly linked list with hash map for O(1) operations
  - Thread-safe implementation patterns
  - Memory management and object pooling concepts

### Design Patterns
- **patterns/observer** shows event-driven architecture:
  - Publisher-subscriber pattern implementation
  - Concurrent event dispatching with goroutines
  - Type-safe event handling and observer management

### Network Programming
- **networking/tcp-echo** demonstrates TCP server patterns:
  - Concurrent client handling with goroutines
  - Protocol design for command handling
  - Connection lifecycle management
  - Interactive and programmatic client modes

### Performance Analysis
- **ctxmemgrowth** provides memory overhead analysis for context usage
- **performance/profiling** demonstrates optimization techniques:
  - Memory allocation patterns and object pooling
  - CPU-intensive task parallelization
  - Algorithm comparison and benchmarking
  - Performance tracking with runtime metrics
- Uses `runtime.ReadMemStats()` for precise memory measurements
- Demonstrates Go runtime introspection patterns
