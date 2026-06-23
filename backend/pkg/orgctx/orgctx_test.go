package orgctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOrgID_RoundTrip(t *testing.T) {
	ctx := WithOrgID(context.Background(), "org-123")
	id, ok := OrgID(ctx)
	assert.True(t, ok)
	assert.Equal(t, "org-123", id)
}

func TestOrgID_Missing(t *testing.T) {
	id, ok := OrgID(context.Background())
	assert.False(t, ok)
	assert.Empty(t, id)
}

func TestOrgID_Empty(t *testing.T) {
	ctx := WithOrgID(context.Background(), "")
	id, ok := OrgID(ctx)
	assert.False(t, ok)
	assert.Empty(t, id)
}
