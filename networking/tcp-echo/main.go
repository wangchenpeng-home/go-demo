package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// TCPServer represents a TCP echo server
type TCPServer struct {
	address  string
	listener net.Listener
	clients  map[net.Conn]string
	mutex    sync.RWMutex
	shutdown chan bool
}

// NewTCPServer creates a new TCP server
func NewTCPServer(address string) *TCPServer {
	return &TCPServer{
		address:  address,
		clients:  make(map[net.Conn]string),
		shutdown: make(chan bool),
	}
}

// Start starts the TCP server
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server on %s: %v", s.address, err)
	}
	
	s.listener = listener
	fmt.Printf("🚀 TCP Echo Server started on %s\n", s.address)
	
	// Start accepting connections in a goroutine
	go s.acceptConnections()
	
	return nil
}

// acceptConnections accepts new client connections
func (s *TCPServer) acceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			// Set a timeout for Accept to allow checking shutdown channel
			if tcpListener, ok := s.listener.(*net.TCPListener); ok {
				tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
			}
			
			conn, err := s.listener.Accept()
			if err != nil {
				// Check if it's a timeout error
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			
			// Handle new client connection
			clientID := fmt.Sprintf("client-%s", conn.RemoteAddr().String())
			s.mutex.Lock()
			s.clients[conn] = clientID
			s.mutex.Unlock()
			
			fmt.Printf("📞 New client connected: %s\n", clientID)
			
			// Handle client in a separate goroutine
			go s.handleClient(conn, clientID)
		}
	}
}

// handleClient handles communication with a single client
func (s *TCPServer) handleClient(conn net.Conn, clientID string) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, conn)
		s.mutex.Unlock()
		conn.Close()
		fmt.Printf("👋 Client disconnected: %s\n", clientID)
	}()
	
	// Send welcome message
	welcome := fmt.Sprintf("Welcome to TCP Echo Server! You are %s\n", clientID)
	conn.Write([]byte(welcome))
	
	// Read and echo messages
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}
		
		// Check for special commands
		if message == "/quit" || message == "/exit" {
			conn.Write([]byte("Goodbye!\n"))
			return
		}
		
		if message == "/clients" {
			s.sendClientList(conn)
			continue
		}
		
		if message == "/time" {
			timeMsg := fmt.Sprintf("Server time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
			conn.Write([]byte(timeMsg))
			continue
		}
		
		// Echo the message back
		echo := fmt.Sprintf("[ECHO] %s: %s\n", time.Now().Format("15:04:05"), message)
		conn.Write([]byte(echo))
		
		fmt.Printf("📨 %s sent: %s\n", clientID, message)
	}
	
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Printf("Error reading from %s: %v", clientID, err)
		}
	}
}

// sendClientList sends the list of connected clients
func (s *TCPServer) sendClientList(conn net.Conn) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	clientList := "Connected clients:\n"
	for _, clientID := range s.clients {
		clientList += fmt.Sprintf("- %s\n", clientID)
	}
	clientList += fmt.Sprintf("Total: %d clients\n", len(s.clients))
	
	conn.Write([]byte(clientList))
}

// Stop stops the TCP server
func (s *TCPServer) Stop() {
	fmt.Println("🛑 Shutting down server...")
	
	// Signal shutdown
	close(s.shutdown)
	
	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}
	
	// Close all client connections
	s.mutex.Lock()
	for conn := range s.clients {
		conn.Write([]byte("Server is shutting down. Goodbye!\n"))
		conn.Close()
	}
	s.mutex.Unlock()
	
	fmt.Println("✅ Server stopped")
}

// TCPClient represents a TCP client
type TCPClient struct {
	serverAddress string
	conn          net.Conn
}

// NewTCPClient creates a new TCP client
func NewTCPClient(serverAddress string) *TCPClient {
	return &TCPClient{
		serverAddress: serverAddress,
	}
}

// Connect connects to the TCP server
func (c *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", c.serverAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to server %s: %v", c.serverAddress, err)
	}
	
	c.conn = conn
	fmt.Printf("✅ Connected to server %s\n", c.serverAddress)
	
	return nil
}

// StartInteractiveSession starts an interactive client session
func (c *TCPClient) StartInteractiveSession() {
	if c.conn == nil {
		fmt.Println("❌ Not connected to server")
		return
	}
	
	// Start reading from server in a goroutine
	go c.readFromServer()
	
	// Read user input and send to server
	fmt.Println("📝 Type messages to send to server. Special commands: /quit, /clients, /time")
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		message := scanner.Text()
		if message == "" {
			continue
		}
		
		// Send message to server
		_, err := c.conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Printf("❌ Error sending message: %v\n", err)
			break
		}
		
		// Check for quit command
		if message == "/quit" || message == "/exit" {
			break
		}
	}
}

// readFromServer reads messages from the server
func (c *TCPClient) readFromServer() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Printf("❌ Error reading from server: %v\n", err)
		}
	}
}

// SendMessage sends a single message to the server
func (c *TCPClient) SendMessage(message string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to server")
	}
	
	// Send message
	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}
	
	// Read response
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response := make([]byte, 1024)
	n, err := c.conn.Read(response)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}
	
	return string(response[:n]), nil
}

// Close closes the client connection
func (c *TCPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		fmt.Println("👋 Disconnected from server")
	}
}

// demonstrateBasicEchoServer shows basic server functionality
func demonstrateBasicEchoServer() {
	fmt.Println("=== Basic TCP Echo Server Demo ===")
	
	// Start server
	server := NewTCPServer("localhost:8080")
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Stop()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Create and connect client
	client := NewTCPClient("localhost:8080")
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	// Send some messages
	messages := []string{
		"Hello, Server!",
		"This is a test message",
		"/time",
		"/clients",
		"Final message",
	}
	
	for _, msg := range messages {
		fmt.Printf("📤 Sending: %s\n", msg)
		response, err := client.SendMessage(msg)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		fmt.Printf("📥 Received: %s", response)
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Println()
}

// demonstrateMultipleClients shows server handling multiple clients
func demonstrateMultipleClients() {
	fmt.Println("=== Multiple Clients Demo ===")
	
	// Start server
	server := NewTCPServer("localhost:8081")
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Stop()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Create multiple clients
	numClients := 3
	clients := make([]*TCPClient, numClients)
	
	// Connect all clients
	for i := 0; i < numClients; i++ {
		clients[i] = NewTCPClient("localhost:8081")
		err := clients[i].Connect()
		if err != nil {
			log.Fatal(err)
		}
		defer clients[i].Close()
		time.Sleep(100 * time.Millisecond) // Stagger connections
	}
	
	// Send messages from each client
	var wg sync.WaitGroup
	for i, client := range clients {
		wg.Add(1)
		go func(clientNum int, c *TCPClient) {
			defer wg.Done()
			
			for j := 0; j < 3; j++ {
				message := fmt.Sprintf("Message %d from client %d", j+1, clientNum+1)
				response, err := c.SendMessage(message)
				if err != nil {
					log.Printf("Client %d error: %v", clientNum+1, err)
					continue
				}
				fmt.Printf("Client %d received: %s", clientNum+1, response)
				time.Sleep(200 * time.Millisecond)
			}
		}(i, client)
	}
	
	wg.Wait()
	
	// Get client list from one client
	response, _ := clients[0].SendMessage("/clients")
	fmt.Printf("Client list:\n%s", response)
	
	fmt.Println()
}

func runInteractiveMode() {
	fmt.Println("=== Interactive Mode ===")
	fmt.Println("Starting server and interactive client...")
	
	// Start server
	server := NewTCPServer("localhost:8082")
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Stop()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Create and connect client
	client := NewTCPClient("localhost:8082")
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	// Start interactive session
	client.StartInteractiveSession()
}

func main() {
	fmt.Println("TCP Echo Server and Client Demo")
	fmt.Println("===============================")
	
	if len(os.Args) > 1 && os.Args[1] == "interactive" {
		runInteractiveMode()
		return
	}
	
	// Run demonstrations
	demonstrateBasicEchoServer()
	time.Sleep(1 * time.Second)
	
	demonstrateMultipleClients()
	
	fmt.Println("✅ TCP demo completed!")
	fmt.Println("💡 Run with 'go run main.go interactive' for interactive mode")
}
