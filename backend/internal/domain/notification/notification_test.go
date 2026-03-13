package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotification_MarkRead(t *testing.T) {
	n := &Notification{
		ID:     "notif-1",
		UserID: "user-1",
		Read:   false,
	}

	assert.False(t, n.Read)
	n.MarkRead()
	assert.True(t, n.Read)
}
