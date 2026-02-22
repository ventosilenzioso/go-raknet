package events

// EventType represents different event types
type EventType int

const (
	EventPlayerConnect EventType = iota
	EventPlayerDisconnect
	EventPlayerSpawn
	EventPlayerDeath
	EventPlayerCommand
	EventPlayerText
	EventPlayerUpdate
	EventVehicleSpawn
	EventVehicleDestroy
)

// Event represents a game event
type Event struct {
	Type      EventType
	PlayerID  uint16
	Data      interface{}
	Timestamp int64
}

// EventHandler is a function that handles events
type EventHandler func(event Event)

// EventManager manages game events
type EventManager struct {
	handlers map[EventType][]EventHandler
}

// NewEventManager creates a new event manager
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Register registers an event handler
func (em *EventManager) Register(eventType EventType, handler EventHandler) {
	em.handlers[eventType] = append(em.handlers[eventType], handler)
}

// Trigger triggers an event
func (em *EventManager) Trigger(event Event) {
	if handlers, exists := em.handlers[event.Type]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
