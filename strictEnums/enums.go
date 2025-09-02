// Package enums demonstrates a strict enum pattern in Go using generics and unexported fields.
// phantom types are used to prevent external construction of enum values.
package enums

import (
	"encoding/json"
	"fmt"
)

type queueType[T any] struct {
	_  [0]T
	id uint8 // unexported - prevents external construction
}

type QueueType = queueType[int]

// Private instances - cannot be modified externally
var (
	fifo       = QueueType{id: 1}
	lifo       = QueueType{id: 2}
	priority   = QueueType{id: 3}
	roundRobin = QueueType{id: 4}
)

// Public constructors - only way to get valid instances
func FIFO() QueueType       { return fifo }
func LIFO() QueueType       { return lifo }
func Priority() QueueType   { return priority }
func RoundRobin() QueueType { return roundRobin }

func (q queueType[int]) String() string {
	switch q.id {
	case 1:
		return "FIFO"
	case 2:
		return "LIFO"
	case 3:
		return "PRIORITY"
	case 4:
		return "ROUND_ROBIN"
	default:
		panic("unreachable") // should never happen
	}
}

// JSON marshaling
func (q queueType[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.String())
}

func (q *queueType[T]) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseQueueType(s)
	if err != nil {
		return err
	}

	*q = any(parsed).(queueType[T])
	return nil
}

func ParseQueueType(s string) (QueueType, error) {
	switch s {
	case "FIFO":
		return FIFO(), nil
	case "LIFO":
		return LIFO(), nil
	case "PRIORITY":
		return Priority(), nil
	case "ROUND_ROBIN":
		return RoundRobin(), nil
	default:
		return QueueType{}, fmt.Errorf("invalid queue type: %s", s)
	}
}

func ProcessQueue(q QueueType) {
	fmt.Printf("Processing %s queue\n", q)
}
