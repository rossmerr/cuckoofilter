package cuckoofilter_test

import (
	"fmt"
	"testing"

	"github.com/rossmerr/cuckoofilter"
)

type test struct {
	sum uint
}

func (s *test) Sum() uint {
	return s.sum
}

func NewTest(s uint) *test {
	return &test{
		sum: s,
	}
}

func TestFilter_Add_Contains(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []uint
		values2 []uint
		want    []uint
	}{
		{
			name:   "Add",
			length: 11,
			values: []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			want:   []uint{20, 50, 53, 75, 100, 67, 3, 36, 39},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](30, 1, 1, 0)

			for _, v := range tt.values {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.want {
				got := filter.Contains(NewTest(v))
				if got != true {
					t.Errorf("Filter.Contains(%v) = %v, want %v", v, got, true)
				}
			}
		})
	}
}

func TestFilter_Matching(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []uint
		values2 []uint
		want    []uint
	}{
		{
			name:   "Add",
			length: 1,
			values: []uint{20, 50, 53, 75},
			want:   []uint{20, 50, 53, 75},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](1, 4, 1, 0)

			for _, v := range tt.values {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.want {
				got := filter.Contains(NewTest(v))
				if got != true {
					t.Errorf("Filter.Contains(%v) = %v, want %v", v, got, true)
				}
			}

			for _, v := range tt.want {

				filter.Remove(NewTest(v))
				got := filter.Contains(NewTest(v))
				if got != false {
					t.Errorf("Filter.Contains() = %v, want %v", got, false)
				}

			}
		})
	}
}

func TestFilter_Add_Rotate_Contains(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []uint
		values2 []uint
		want    []uint
	}{
		{
			name:    "Add",
			length:  11,
			values:  []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			values2: []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39, 6},
			want:    []uint{20, 50, 53, 75, 100, 67, 3, 36, 39, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](30, 1, 1, 0)

			for _, v := range tt.values {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.values2 {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.want {
				got := filter.Contains(NewTest(v))
				if got != true {
					t.Errorf("Filter.Contains() = %v, want %v", got, true)
				}
			}
		})
	}
}

func TestFilter_Add_Delete_Contains(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []uint
		values2 []uint
		want    []uint
	}{
		{
			name:    "Add",
			length:  11,
			values:  []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			values2: []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39, 6},
			want:    []uint{20, 50, 53, 75, 100, 67, 3, 36, 39, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](30, 4, 1, 11)

			for _, v := range tt.values {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.values2 {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for i, v := range tt.want {
				check := tt.want[i+1:]

				filter.Remove(NewTest(v))
				got := filter.Contains(NewTest(v))
				if got != false {
					t.Errorf("Filter.Contains(%v) = %v, want %v", v, got, false)
				}
				for _, v := range check {
					got := filter.Contains(NewTest(v))
					if got != true {
						t.Errorf("Filter.Contains() = %v, want %v", got, true)
					}
				}
			}
		})
	}
}

func TestFilterFalsePositiveRate_Add_Contains(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []uint
		values2 []uint
		want    []uint
	}{
		{
			name:    "Add",
			length:  11,
			values:  []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			values2: []uint{20, 50, 53, 75, 100, 67, 105, 3, 36, 39, 6},
			want:    []uint{20, 50, 53, 75, 100, 67, 3, 36, 39, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilterFalsePositiveRate[*test](uint(tt.length), 0.1)

			for _, v := range tt.values {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.values2 {
				err := filter.Add(NewTest(v))
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, v := range tt.want {
				got := filter.Contains(NewTest(v))
				if got != true {
					t.Errorf("Filter.Contains() = %v, want %v", got, true)
				}
			}
		})
	}
}
