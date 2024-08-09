package util

import "testing"

func TestCheckModel(t *testing.T) {
	type args struct {
		Name  string  `json:"name"`
		Age   int     `json:"age"`
		Score float64 `json:"score"`
	}

	type testData struct {
		s     []string
		valid bool
	}

	tests := []testData{
		{
			s:     []string{"name", "age", "score"},
			valid: false,
		},
		{
			s:     []string{"hello", "10", "100.0"},
			valid: true,
		},
		{
			s:     []string{"hello", "10.1", "100.0"},
			valid: false,
		},
		{
			s:     []string{"hello", "10", "100"},
			valid: true,
		},
	}

	wrapErr2Valid := func(err error) bool {
		return err == nil
	}

	a := &args{}
	for _, tt := range tests {
		t.Run("TestCheckModel", func(t *testing.T) {
			if got := CheckModelOrder(a, tt.s); wrapErr2Valid(got) != tt.valid {
				t.Errorf("CheckModel() = %v, want %v", got, tt.valid)
			}
		})
	}
}
