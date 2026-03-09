package notification

import "time"

type Channel string

const (
	ChannelEmail  Channel = "email"
	ChannelInApp  Channel = "in_app"
)

type Notification struct {
	ID        string
	UserID    string
	Title     string
	Body      string
	Channel   Channel
	Read      bool
	CreatedAt time.Time
}

func (n *Notification) MarkRead() {
	n.Read = true
}
