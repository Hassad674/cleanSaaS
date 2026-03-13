package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	svc := NewService(nil, nil)
	assert.NotNil(t, svc, "NewService should return a non-nil service")
}
