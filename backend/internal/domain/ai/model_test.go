package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		model Model
		valid bool
	}{
		{"Claude is valid", ModelClaude, true},
		{"GPT is valid", ModelGPT, true},
		{"Gemini is valid", ModelGemini, true},
		{"Empty string is invalid", Model(""), false},
		{"Random string is invalid", Model("llama"), false},
		{"Case sensitive", Model("Claude"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.model.IsValid())
		})
	}
}

func TestModel_Constants(t *testing.T) {
	assert.Equal(t, Model("claude"), ModelClaude)
	assert.Equal(t, Model("gpt"), ModelGPT)
	assert.Equal(t, Model("gemini"), ModelGemini)
}
