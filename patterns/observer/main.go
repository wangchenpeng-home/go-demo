package main

import (
	"fmt"
	"sync"
	"time"
)

// Observer interface defines the contract for observers
type Observer interface {
	Update(event Event)
	GetID() string
}

// Subject interface defines the contract for subjects
type Subject interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Notify(event Event)
}

// Event represents an event in the system
type Event struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
	Source    string
}

// EventDispatcher manages observers and handles event distribution
type EventDispatcher struct {
	observers map[string][]Observer
	mutex     sync.RWMutex
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		observers: make(map[string][]Observer),
	}
}

// Subscribe adds an observer for a specific event type
func (ed *EventDispatcher) Subscribe(eventType string, observer Observer) {
	ed.mutex.Lock()
	defer ed.mutex.Unlock()
	
	ed.observers[eventType] = append(ed.observers[eventType], observer)
	fmt.Printf("Observer %s subscribed to event type: %s\n", observer.GetID(), eventType)
}

// Unsubscribe removes an observer for a specific event type
func (ed *EventDispatcher) Unsubscribe(eventType string, observer Observer) {
	ed.mutex.Lock()
	defer ed.mutex.Unlock()
	
	observers := ed.observers[eventType]
	for i, obs := range observers {
		if obs.GetID() == observer.GetID() {
			ed.observers[eventType] = append(observers[:i], observers[i+1:]...)
			fmt.Printf("Observer %s unsubscribed from event type: %s\n", observer.GetID(), eventType)
			return
		}
	}
}

// Notify sends an event to all subscribed observers
func (ed *EventDispatcher) Notify(event Event) {
	ed.mutex.RLock()
	observers := ed.observers[event.Type]
	ed.mutex.RUnlock()
	
	fmt.Printf("Broadcasting event: %s from %s\n", event.Type, event.Source)
	
	// Notify all observers concurrently
	var wg sync.WaitGroup
	for _, observer := range observers {
		wg.Add(1)
		go func(obs Observer) {
			defer wg.Done()
			obs.Update(event)
		}(observer)
	}
	wg.Wait()
}

// StockPrice represents a stock price subject
type StockPrice struct {
	Symbol     string
	Price      float64
	Change     float64
	dispatcher *EventDispatcher
}

// NewStockPrice creates a new stock price subject
func NewStockPrice(symbol string, dispatcher *EventDispatcher) *StockPrice {
	return &StockPrice{
		Symbol:     symbol,
		dispatcher: dispatcher,
	}
}

// SetPrice updates the stock price and notifies observers
func (sp *StockPrice) SetPrice(newPrice float64) {
	oldPrice := sp.Price
	sp.Price = newPrice
	sp.Change = newPrice - oldPrice
	
	event := Event{
		Type:      "price_update",
		Data:      sp,
		Timestamp: time.Now(),
		Source:    fmt.Sprintf("StockPrice-%s", sp.Symbol),
	}
	
	sp.dispatcher.Notify(event)
}

// EmailNotifier implements Observer for email notifications
type EmailNotifier struct {
	ID    string
	Email string
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(id, email string) *EmailNotifier {
	return &EmailNotifier{
		ID:    id,
		Email: email,
	}
}

// Update handles incoming events
func (en *EmailNotifier) Update(event Event) {
	switch event.Type {
	case "price_update":
		if stock, ok := event.Data.(*StockPrice); ok {
			fmt.Printf("📧 EMAIL to %s: %s price changed to $%.2f (change: %+.2f) at %s\n",
				en.Email, stock.Symbol, stock.Price, stock.Change, event.Timestamp.Format("15:04:05"))
		}
	case "system_alert":
		fmt.Printf("📧 EMAIL ALERT to %s: %s at %s\n",
			en.Email, event.Data, event.Timestamp.Format("15:04:05"))
	}
	
	// Simulate email sending delay
	time.Sleep(50 * time.Millisecond)
}

// GetID returns the observer ID
func (en *EmailNotifier) GetID() string {
	return en.ID
}

// SMSNotifier implements Observer for SMS notifications
type SMSNotifier struct {
	ID    string
	Phone string
}

// NewSMSNotifier creates a new SMS notifier
func NewSMSNotifier(id, phone string) *SMSNotifier {
	return &SMSNotifier{
		ID:    id,
		Phone: phone,
	}
}

// Update handles incoming events
func (sn *SMSNotifier) Update(event Event) {
	switch event.Type {
	case "price_update":
		if stock, ok := event.Data.(*StockPrice); ok {
			if stock.Change > 5.0 || stock.Change < -5.0 {
				fmt.Printf("📱 SMS to %s: ALERT! %s significant change: $%.2f (%+.2f) at %s\n",
					sn.Phone, stock.Symbol, stock.Price, stock.Change, event.Timestamp.Format("15:04:05"))
			}
		}
	case "system_alert":
		fmt.Printf("📱 SMS ALERT to %s: %s at %s\n",
			sn.Phone, event.Data, event.Timestamp.Format("15:04:05"))
	}
	
	// Simulate SMS sending delay
	time.Sleep(30 * time.Millisecond)
}

// GetID returns the observer ID
func (sn *SMSNotifier) GetID() string {
	return sn.ID
}

// LoggingObserver implements Observer for logging events
type LoggingObserver struct {
	ID       string
	LogLevel string
}

// NewLoggingObserver creates a new logging observer
func NewLoggingObserver(id, logLevel string) *LoggingObserver {
	return &LoggingObserver{
		ID:       id,
		LogLevel: logLevel,
	}
}

// Update handles incoming events
func (lo *LoggingObserver) Update(event Event) {
	switch event.Type {
	case "price_update":
		if stock, ok := event.Data.(*StockPrice); ok {
			fmt.Printf("📝 LOG [%s]: Price update - %s: $%.2f (change: %+.2f) at %s\n",
				lo.LogLevel, stock.Symbol, stock.Price, stock.Change, event.Timestamp.Format("15:04:05"))
		}
	case "system_alert":
		fmt.Printf("📝 LOG [%s]: System Alert - %s at %s\n",
			lo.LogLevel, event.Data, event.Timestamp.Format("15:04:05"))
	}
}

// GetID returns the observer ID
func (lo *LoggingObserver) GetID() string {
	return lo.ID
}

// DatabaseObserver implements Observer for database operations
type DatabaseObserver struct {
	ID string
}

// NewDatabaseObserver creates a new database observer
func NewDatabaseObserver(id string) *DatabaseObserver {
	return &DatabaseObserver{ID: id}
}

// Update handles incoming events
func (do *DatabaseObserver) Update(event Event) {
	switch event.Type {
	case "price_update":
		if stock, ok := event.Data.(*StockPrice); ok {
			fmt.Printf("💾 DB: Saved price record - %s: $%.2f at %s\n",
				stock.Symbol, stock.Price, event.Timestamp.Format("15:04:05"))
		}
	}
	
	// Simulate database write delay
	time.Sleep(80 * time.Millisecond)
}

// GetID returns the observer ID
func (do *DatabaseObserver) GetID() string {
	return do.ID
}

func demonstrateBasicObserverPattern() {
	fmt.Println("=== Basic Observer Pattern Demo ===")
	
	// Create event dispatcher
	dispatcher := NewEventDispatcher()
	
	// Create observers
	emailNotifier := NewEmailNotifier("email1", "trader@example.com")
	smsNotifier := NewSMSNotifier("sms1", "+1234567890")
	logger := NewLoggingObserver("logger1", "INFO")
	
	// Subscribe observers to events
	dispatcher.Subscribe("price_update", emailNotifier)
	dispatcher.Subscribe("price_update", smsNotifier)
	dispatcher.Subscribe("price_update", logger)
	dispatcher.Subscribe("system_alert", emailNotifier)
	dispatcher.Subscribe("system_alert", smsNotifier)
	dispatcher.Subscribe("system_alert", logger)
	
	// Create stock and trigger some price updates
	stock := NewStockPrice("AAPL", dispatcher)
	stock.SetPrice(150.0)
	time.Sleep(100 * time.Millisecond)
	
	stock.SetPrice(155.5)
	time.Sleep(100 * time.Millisecond)
	
	// Trigger a system alert
	alertEvent := Event{
		Type:      "system_alert",
		Data:      "Trading system maintenance scheduled",
		Timestamp: time.Now(),
		Source:    "SystemManager",
	}
	dispatcher.Notify(alertEvent)
	
	fmt.Println()
}

func demonstrateAdvancedScenario() {
	fmt.Println("=== Advanced Scenario - Stock Trading System ===")
	
	dispatcher := NewEventDispatcher()
	
	// Create multiple observers
	observers := []Observer{
		NewEmailNotifier("email1", "trader1@example.com"),
		NewEmailNotifier("email2", "trader2@example.com"),
		NewSMSNotifier("sms1", "+1234567890"),
		NewSMSNotifier("sms2", "+0987654321"),
		NewLoggingObserver("logger1", "INFO"),
		NewDatabaseObserver("db1"),
	}
	
	// Subscribe all observers to price updates
	for _, observer := range observers {
		dispatcher.Subscribe("price_update", observer)
	}
	
	// Create multiple stocks
	stocks := []*StockPrice{
		NewStockPrice("AAPL", dispatcher),
		NewStockPrice("GOOGL", dispatcher),
		NewStockPrice("MSFT", dispatcher),
	}
	
	// Simulate market activity
	fmt.Println("Simulating market activity...")
	prices := [][]float64{
		{150.0, 152.5, 148.0, 155.0},
		{2500.0, 2520.0, 2480.0, 2510.0},
		{300.0, 305.0, 298.0, 310.0},
	}
	
	for i, priceList := range prices {
		stock := stocks[i]
		for _, price := range priceList {
			stock.SetPrice(price)
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	// Demonstrate unsubscribing
	fmt.Println("\nUnsubscribing SMS notifier...")
	dispatcher.Unsubscribe("price_update", observers[2]) // Remove sms1
	
	// Continue with more price updates
	stocks[0].SetPrice(160.0)
	time.Sleep(100 * time.Millisecond)
	
	fmt.Println()
}

func main() {
	fmt.Println("Observer Pattern Implementation Demo")
	fmt.Println("===================================")
	
	demonstrateBasicObserverPattern()
	time.Sleep(500 * time.Millisecond)
	
	demonstrateAdvancedScenario()
	
	fmt.Println("Observer pattern demo completed!")
}
