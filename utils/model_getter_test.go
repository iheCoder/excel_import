package util

import "testing"

type fieldGetterTest1 struct {
	A string
	B int
	C float64
}

func TestGetFieldString(t *testing.T) {
	type testData struct {
		st       any
		expected []string
	}

	tests := []testData{
		{
			st: &fieldGetterTest1{
				A: "hello",
				B: 10,
				C: 100.0,
			},
			expected: []string{"hello", "10", "100"},
		},
		{
			st: &fieldGetterTest1{
				A: "world",
				B: 20,
				C: 200.01122124562,
			},
			expected: []string{"world", "20", "200.01122124562"},
		},
	}

	numFields := 3
	for _, tt := range tests {
		t.Run("TestGetFieldString", func(t *testing.T) {
			for i := 0; i < numFields; i++ {
				if got, _ := GetFieldString(tt.st, i); got != tt.expected[i] {
					t.Errorf("GetFieldString() = %v, want %v", got, tt.expected[i])
				}
			}
		})
	}
}
