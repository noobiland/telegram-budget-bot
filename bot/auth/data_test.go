package auth

import (
	"reflect"
	"testing"
)

// InitUsers test data preparation
type initUsersTestCase = struct {
	name           string
	mockData       map[int]string
	expectedResult map[int]string
	expectError    bool
}

var initUsersTestCases = []initUsersTestCase{
	{
		name:           "empty user map",
		mockData:       map[int]string{},
		expectedResult: map[int]string{},
		expectError:    false,
	},
	{
		name: "single user",
		mockData: map[int]string{
			1: "Alice",
		},
		expectedResult: map[int]string{
			1: "Alice",
		},
	},
	{
		name:           "nil user map",
		mockData:       nil,
		expectedResult: map[int]string{},
		expectError:    true,
	},
}

func TestInitUsers(t *testing.T) {

	for _, tc := range initUsersTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitUsers(func() map[int]string {
				return tc.mockData
			})

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(Ids, tc.expectedResult) {
					t.Errorf("expected %v, got %v", tc.expectedResult, Ids)
				}
			}
		})
	}
}
