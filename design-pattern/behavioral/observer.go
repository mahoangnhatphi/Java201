package behavioral

import (
	"fmt"
	"sync"
)

// EventType represents the type of event
type EventType string

const (
	UserCreated EventType = "user.created"
	UserUpdated EventType = "user.updated"
	UserDeleted EventType = "user.deleted"
	OrderPlaced EventType = "order.placed"
)

// Event represents an event in the system
type Event struct {
	Type    EventType
	Payload interface{}
}

// Observer defines the interface for observers
type Observer interface {
	OnNotify(event Event)
}

// Subject defines the interface for subjects that can be observed
type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	Notify(event Event)
}

// EventBus is a concrete implementation of Subject
type EventBus struct {
	observers map[EventType][]Observer
	mu        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[EventType][]Observer),
	}
}

func (e *EventBus) Register(observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Register for all events by default
	for eventType := range observerTypes(observer) {
		e.observers[eventType] = append(e.observers[eventType], observer)
	}
}

func (e *EventBus) Deregister(observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for eventType, observers := range e.observers {
		filtered := observers[:0]
		for _, obs := range observers {
			if obs != observer {
				filtered = append(filtered, obs)
			}
		}
		e.observers[eventType] = filtered
	}
}

func (e *EventBus) Notify(event Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	observers, ok := e.observers[event.Type]
	if !ok {
		return
	}

	for _, observer := range observers {
		observer.OnNotify(event)
	}
}

// observerTypes returns the event types an observer is interested in
// This is a simplified implementation - in real systems you'd have explicit registration
func observerTypes(observer Observer) map[EventType]bool {
	return map[EventType]bool{
		UserCreated: true,
		UserUpdated: true,
		UserDeleted: true,
		OrderPlaced: true,
	}
}

// EmailNotifier is an observer that sends emails
type EmailNotifier struct {
	emailService string
}

func NewEmailNotifier(emailService string) *EmailNotifier {
	return &EmailNotifier{emailService: emailService}
}

func (e *EmailNotifier) OnNotify(event Event) {
	fmt.Printf("[EmailNotifier] Sending email via %s for event: %s\n", e.emailService, event.Type)
}

// LoggerObserver is an observer that logs events
type LoggerObserver struct{}

func NewLoggerObserver() *LoggerObserver {
	return &LoggerObserver{}
}

func (l *LoggerObserver) OnNotify(event Event) {
	fmt.Printf("[Logger] Event received: %s with payload: %v\n", event.Type, event.Payload)
}

// AnalyticsObserver is an observer that tracks analytics
type AnalyticsObserver struct {
	trackingID string
}

func NewAnalyticsObserver(trackingID string) *AnalyticsObserver {
	return &AnalyticsObserver{trackingID: trackingID}
}

func (a *AnalyticsObserver) OnNotify(event Event) {
	fmt.Printf("[Analytics] Tracking event %s with ID %s\n", event.Type, a.trackingID)
}

// ObserverExampleUsage demonstrates the Observer pattern
func ObserverExampleUsage() {
	eventBus := NewEventBus()

	emailNotifier := NewEmailNotifier("SendGrid")
	logger := NewLoggerObserver()
	analytics := NewAnalyticsObserver("GA-123456")

	eventBus.Register(emailNotifier)
	eventBus.Register(logger)
	eventBus.Register(analytics)

	eventBus.Notify(Event{
		Type:    UserCreated,
		Payload: map[string]string{"user_id": "123", "email": "user@example.com"},
	})

	eventBus.Notify(Event{
		Type:    OrderPlaced,
		Payload: map[string]interface{}{"order_id": "456", "total": 99.99},
	})
}
