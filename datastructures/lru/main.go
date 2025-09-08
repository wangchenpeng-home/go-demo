package main

import (
	"fmt"
	"time"
)

// Node represents a node in the doubly linked list
type Node struct {
	Key, Value interface{}
	Prev, Next *Node
}

// LRUCache represents an LRU cache
type LRUCache struct {
	capacity int
	cache    map[interface{}]*Node
	head     *Node // Most recently used
	tail     *Node // Least recently used
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	
	// Create dummy head and tail nodes
	head := &Node{}
	tail := &Node{}
	head.Next = tail
	tail.Prev = head
	
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[interface{}]*Node),
		head:     head,
		tail:     tail,
	}
}

// addNode adds a node right after head
func (lru *LRUCache) addNode(node *Node) {
	node.Prev = lru.head
	node.Next = lru.head.Next
	
	lru.head.Next.Prev = node
	lru.head.Next = node
}

// removeNode removes an existing node from the linked list
func (lru *LRUCache) removeNode(node *Node) {
	prevNode := node.Prev
	nextNode := node.Next
	
	prevNode.Next = nextNode
	nextNode.Prev = prevNode
}

// moveToHead moves certain node to the head
func (lru *LRUCache) moveToHead(node *Node) {
	lru.removeNode(node)
	lru.addNode(node)
}

// popTail removes the last node (LRU node)
func (lru *LRUCache) popTail() *Node {
	lastNode := lru.tail.Prev
	lru.removeNode(lastNode)
	return lastNode
}

// Get retrieves a value by key
func (lru *LRUCache) Get(key interface{}) (interface{}, bool) {
	node, exists := lru.cache[key]
	if !exists {
		return nil, false
	}
	
	// Move the accessed node to the head (mark as recently used)
	lru.moveToHead(node)
	
	return node.Value, true
}

// Put sets a key-value pair in the cache
func (lru *LRUCache) Put(key, value interface{}) {
	node, exists := lru.cache[key]
	
	if exists {
		// Update the value and move to head
		node.Value = value
		lru.moveToHead(node)
		return
	}
	
	// Create new node
	newNode := &Node{
		Key:   key,
		Value: value,
	}
	
	if len(lru.cache) >= lru.capacity {
		// Remove LRU node
		tail := lru.popTail()
		delete(lru.cache, tail.Key)
	}
	
	// Add new node
	lru.addNode(newNode)
	lru.cache[key] = newNode
}

// Size returns the current size of the cache
func (lru *LRUCache) Size() int {
	return len(lru.cache)
}

// Capacity returns the capacity of the cache
func (lru *LRUCache) Capacity() int {
	return lru.capacity
}

// Keys returns all keys in order from most recent to least recent
func (lru *LRUCache) Keys() []interface{} {
	keys := make([]interface{}, 0, len(lru.cache))
	current := lru.head.Next
	for current != lru.tail {
		keys = append(keys, current.Key)
		current = current.Next
	}
	return keys
}

// Clear removes all entries from the cache
func (lru *LRUCache) Clear() {
	lru.cache = make(map[interface{}]*Node)
	lru.head.Next = lru.tail
	lru.tail.Prev = lru.head
}

// String returns a string representation of the cache
func (lru *LRUCache) String() string {
	keys := lru.Keys()
	result := fmt.Sprintf("LRUCache{capacity: %d, size: %d, keys: [", lru.capacity, len(lru.cache))
	for i, key := range keys {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v", key)
	}
	result += "]}"
	return result
}

// Demo function to show LRU cache in action
func demonstrateBasicOperations() {
	fmt.Println("=== Basic LRU Cache Operations ===")
	cache := NewLRUCache(3)
	
	// Test Put operations
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)
	fmt.Printf("After adding a, b, c: %s\n", cache.String())
	
	// Test Get operation
	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %v\n", val)
	}
	fmt.Printf("After accessing 'a': %s\n", cache.String())
	
	// Test eviction
	cache.Put("d", 4) // This should evict 'b'
	fmt.Printf("After adding 'd' (should evict 'b'): %s\n", cache.String())
	
	// Test accessing non-existent key
	if _, ok := cache.Get("b"); !ok {
		fmt.Println("Key 'b' not found (correctly evicted)")
	}
	
	fmt.Println()
}

// Demo function to show cache performance characteristics
func demonstratePerformance() {
	fmt.Println("=== Performance Demonstration ===")
	cache := NewLRUCache(1000)
	
	// Fill cache with data
	start := time.Now()
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key%d", i), i*i)
	}
	fillTime := time.Since(start)
	fmt.Printf("Time to fill cache with 1000 items: %v\n", fillTime)
	
	// Test random access performance
	start = time.Now()
	hits := 0
	accesses := 10000
	for i := 0; i < accesses; i++ {
		key := fmt.Sprintf("key%d", i%1200) // Some keys won't exist
		if _, ok := cache.Get(key); ok {
			hits++
		}
	}
	accessTime := time.Since(start)
	fmt.Printf("Time for %d random accesses: %v\n", accesses, accessTime)
	fmt.Printf("Hit rate: %.2f%%\n", float64(hits)/float64(accesses)*100)
	
	fmt.Println()
}

// Demo function to show LRU eviction policy
func demonstrateEvictionPolicy() {
	fmt.Println("=== LRU Eviction Policy Demonstration ===")
	cache := NewLRUCache(4)
	
	// Fill cache
	for i := 1; i <= 4; i++ {
		cache.Put(i, fmt.Sprintf("value%d", i))
	}
	fmt.Printf("Initial state: %s\n", cache.String())
	
	// Access some items to change their order
	cache.Get(2)
	fmt.Printf("After accessing key 2: %s\n", cache.String())
	
	cache.Get(4)
	fmt.Printf("After accessing key 4: %s\n", cache.String())
	
	// Add new item - should evict key 1 (least recently used)
	cache.Put(5, "value5")
	fmt.Printf("After adding key 5 (should evict key 1): %s\n", cache.String())
	
	// Add another item - should evict key 3
	cache.Put(6, "value6")
	fmt.Printf("After adding key 6 (should evict key 3): %s\n", cache.String())
	
	fmt.Println()
}

// Demo function showing practical use case - web page caching
func demonstrateWebPageCache() {
	fmt.Println("=== Web Page Cache Use Case ===")
	pageCache := NewLRUCache(5)
	
	// Simulate web page requests
	pages := []string{
		"/home", "/about", "/products", "/contact", "/blog",
		"/home", "/news", "/about", "/services", "/home",
	}
	
	for _, page := range pages {
		if content, exists := pageCache.Get(page); exists {
			fmt.Printf("Cache HIT: %s -> %s\n", page, content)
		} else {
			// Simulate loading page content
			content := fmt.Sprintf("Content of %s page", page)
			pageCache.Put(page, content)
			fmt.Printf("Cache MISS: Loaded %s\n", page)
		}
		fmt.Printf("  Current cache: %s\n", pageCache.String())
	}
	
	fmt.Println()
}

func main() {
	fmt.Println("LRU Cache Implementation Demo")
	fmt.Println("============================")
	
	demonstrateBasicOperations()
	demonstratePerformance()
	demonstrateEvictionPolicy()
	demonstrateWebPageCache()
	
	fmt.Println("Demo completed!")
}
