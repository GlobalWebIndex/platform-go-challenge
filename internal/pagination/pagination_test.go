package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePaginationParams_Defaults(t *testing.T) {
	p := ParsePaginationParams("", "")
	assert.Equal(t, "", p.Cursor)
	assert.Equal(t, DefaultLimit, p.Limit)
}

func TestParsePaginationParams_ParsesLimit(t *testing.T) {
	p := ParsePaginationParams("abc", "10")
	assert.Equal(t, "abc", p.Cursor)
	assert.Equal(t, 10, p.Limit)
}

func TestParsePaginationParams_InvalidLimitStringFallsBackToDefault(t *testing.T) {
	p := ParsePaginationParams("abc", "not-a-number")
	assert.Equal(t, "abc", p.Cursor)
	assert.Equal(t, DefaultLimit, p.Limit)
}

func TestParsePaginationParams_EnforcesMinLimit(t *testing.T) {
	p := ParsePaginationParams("", "0")
	assert.Equal(t, DefaultLimit, p.Limit)

	p2 := ParsePaginationParams("", "-5")
	assert.Equal(t, DefaultLimit, p2.Limit)
}

func TestParsePaginationParams_EnforcesMaxLimit(t *testing.T) {
	p := ParsePaginationParams("", "200")
	assert.Equal(t, MaxLimit, p.Limit)
}

func TestNewPaginatedResponse(t *testing.T) {
	data := []int{1, 2, 3}
	resp := NewPaginatedResponse(data, "next", true, 20)

	assert.Equal(t, data, resp.Data)

	assert.Equal(t, "next", resp.Meta.NextCursor)
	assert.True(t, resp.Meta.HasMore)
	assert.Equal(t, 20, resp.Meta.Limit)
	assert.Equal(t, 0, resp.Meta.Total) // not set yet
}

func TestPaginatedResponse_WithTotal(t *testing.T) {
	resp := NewPaginatedResponse([]string{"a"}, "", false, 10)
	out := resp.WithTotal(123)

	assert.Same(t, resp, out)
	assert.Equal(t, 123, resp.Meta.Total)
}
