package api

import (
	"strings"
	"testing"
)

func TestBuildQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "Empty params",
			params:   map[string]interface{}{},
			expected: "",
		},
		{
			name: "Single param",
			params: map[string]interface{}{
				"key": "value",
			},
			expected: "key=value",
		},
		{
			name: "Multiple params",
			params: map[string]interface{}{
				"limit":      10,
				"includeAll": true,
			},
			expected: "includeAll=true&limit=10", // Note: order might vary due to map iteration
		},
		{
			name: "Params with nil values should be skipped",
			params: map[string]interface{}{
				"key1": nil,
				"key2": "value2",
			},
			expected: "key2=value2",
		},
		{
			name: "Empty string values should be skipped",
			params: map[string]interface{}{
				"key1": "",
				"key2": "value2",
			},
			expected: "key2=value2",
		},
		{
			name: "Empty slice should be skipped",
			params: map[string]interface{}{
				"emptySlice": []int{},
				"key2":       "value2",
			},
			expected: "key2=value2",
		},
		{
			name: "Non-empty slice should be included",
			params: map[string]interface{}{
				"slice": []int{1, 2, 3},
				"key":   "value",
			},
			expected: "key=value&slice=[1 2 3]", // Note: order might vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildQueryParams(tt.params)

			// For tests with multiple params, we need to check if all expected params are present
			// since map iteration order is not guaranteed
			if len(tt.params) <= 1 {
				if result != tt.expected {
					t.Errorf("buildQueryParams() = %v, want %v", result, tt.expected)
				}
			} else {
				// For multiple params, check that the result contains all expected key-value pairs
				expectedPairs := strings.Split(tt.expected, "&")
				resultPairs := strings.Split(result, "&")

				if len(expectedPairs) != len(resultPairs) {
					t.Errorf("buildQueryParams() = %v, want %v", result, tt.expected)
					return
				}

				// Check if all expected pairs are present (order doesn't matter)
				for _, expectedPair := range expectedPairs {
					found := false
					for _, resultPair := range resultPairs {
						if expectedPair == resultPair {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("buildQueryParams() missing expected pair %v in result %v", expectedPair, result)
					}
				}
			}
		})
	}
}

func BenchmarkBuildQueryParams(b *testing.B) {
	params := map[string]interface{}{
		"limit":      100,
		"cursor":     "abc123",
		"includeAll": true,
		"filter":     "status:active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildQueryParams(params)
	}
}
