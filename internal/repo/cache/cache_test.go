package cache

import (
	"testing"
)

func TestGet(t *testing.T) {
	cache := NewLRUCache[string, string](2)

	tests := []struct {
		name      string
		operation func()
		key       string
		expected  string
		found     bool
	}{
		{
			name:      "Test empty cache",
			operation: func() {},
			key:       "key1",
			expected:  "",
			found:     false,
		},
		{
			name: "Test key1 addition",
			operation: func() {
				cache.Put("key1", "content1")
			},
			key:      "key1",
			expected: "content1",
			found:    true,
		},
		{
			name: "Test key2 addition",
			operation: func() {
				cache.Put("key2", "content2")
			},
			key:      "key2",
			expected: "content2",
			found:    true,
		},
		{
			name: "Test key3 addition and cache bounds",
			operation: func() {
				cache.Put("key3", "content3")
			},
			key:      "key1",
			expected: "",
			found:    false,
		},
		{
			name:      "Test key2 after adding key3",
			operation: func() {},
			key:       "key2",
			expected:  "content2",
			found:     true,
		},
		{
			name:      "Test key3 after adding key3",
			operation: func() {},
			key:       "key3",
			expected:  "content3",
			found:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.operation()
			value, found := cache.Get(tt.key)
			if found != tt.found {
				t.Errorf("Expected found=%v, received=%v", tt.found, found)
			}
			if value != tt.expected {
				t.Errorf("Expected value=%v, received=%v", tt.expected, value)
			}
		})
	}
}
