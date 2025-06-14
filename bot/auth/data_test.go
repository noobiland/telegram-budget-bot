package auth

import (
	"reflect"
	"testing"
)

func TestInitUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockData       map[int]string
		expectedResult map[int]string
		expectError    bool
	}{
		{
			name:           "empty user map",
			mockData:       map[int]string{},
			expectedResult: map[int]string{},
			expectError:    false,
		},
		{
			name: "single user",
			mockData: map[int]string{
				1: "UserA",
			},
			expectedResult: map[int]string{
				1: "UserA",
			},
		},
		{
			name:           "nil user map",
			mockData:       nil,
			expectedResult: map[int]string{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitUsers(func() map[int]string {
				return tt.mockData
			})

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(Ids, tt.expectedResult) {
					t.Errorf("expected %v, got %v", tt.expectedResult, Ids)
				}
			}
		})
	}
}
