package main

import "testing"

func Test_openAPI(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test-openAPI"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openAPI()
		})
	}
}
