package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromRequest_Defaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	page := FromRequest(req)

	assert.Equal(t, 0, page.Offset)
	assert.Equal(t, 20, page.Limit)
}

func TestFromRequest_CustomValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?offset=10&limit=50", nil)
	page := FromRequest(req)

	assert.Equal(t, 10, page.Offset)
	assert.Equal(t, 50, page.Limit)
}

func TestFromRequest_NegativeOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?offset=-5", nil)
	page := FromRequest(req)

	assert.Equal(t, 0, page.Offset, "negative offset should use default")
}

func TestFromRequest_ZeroLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=0", nil)
	page := FromRequest(req)

	assert.Equal(t, 20, page.Limit, "zero limit should use default")
}

func TestFromRequest_NegativeLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=-10", nil)
	page := FromRequest(req)

	assert.Equal(t, 20, page.Limit, "negative limit should use default")
}

func TestFromRequest_OverMaxLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=200", nil)
	page := FromRequest(req)

	assert.Equal(t, 20, page.Limit, "limit over 100 should use default")
}

func TestFromRequest_ExactMaxLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=100", nil)
	page := FromRequest(req)

	assert.Equal(t, 100, page.Limit, "limit of exactly 100 should be accepted")
}

func TestFromRequest_InvalidOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?offset=abc", nil)
	page := FromRequest(req)

	assert.Equal(t, 0, page.Offset, "non-numeric offset should use default")
}

func TestFromRequest_InvalidLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=xyz", nil)
	page := FromRequest(req)

	assert.Equal(t, 20, page.Limit, "non-numeric limit should use default")
}

func TestFromRequest_OnlyOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?offset=30", nil)
	page := FromRequest(req)

	assert.Equal(t, 30, page.Offset)
	assert.Equal(t, 20, page.Limit, "missing limit should use default")
}

func TestFromRequest_OnlyLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items?limit=5", nil)
	page := FromRequest(req)

	assert.Equal(t, 0, page.Offset, "missing offset should use default")
	assert.Equal(t, 5, page.Limit)
}

func TestPage_JSONTags(t *testing.T) {
	// Page struct should have correct JSON tags for serialization
	page := Page{Offset: 10, Limit: 20, Total: 100}
	assert.Equal(t, 10, page.Offset)
	assert.Equal(t, 20, page.Limit)
	assert.Equal(t, 100, page.Total)
}
