package notification

import "time"

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Title     string
	Message   string
	Read      bool
	Data      map[string]interface{}
	CreatedAt time.Time
}

func (n *Notification) MarkRead() {
	n.Read = true
}
