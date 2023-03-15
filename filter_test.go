package cuckoofilter_test

import (
	"fmt"
	"testing"

	"github.com/rossmerr/cuckoofilter"
)

type test struct {
	sum int
}

func (s *test) Sum() int {
	return s.sum
}

func NewTest(s int) *test {
	return &test{
		sum: s,
	}
}

func TestFilter_Add_Contains(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		values  []int
		values2 []int
		want    []int
	}{
		{
			name:    "Add",
			length:  11,
			values:  []int{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			values2: []int{20, 50, 53, 75, 100, 67, 105, 3, 36, 39, 6},
			want:    []int{20, 50, 53, 75, 100, 67, 3, 36, 39, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](2, tt.length)

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
		values  []int
		values2 []int
		want    []int
	}{
		{
			name:    "Add",
			length:  11,
			values:  []int{20, 50, 53, 75, 100, 67, 105, 3, 36, 39},
			values2: []int{20, 50, 53, 75, 100, 67, 105, 3, 36, 39, 6},
			want:    []int{20, 50, 53, 75, 100, 67, 3, 36, 39, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := cuckoofilter.NewFilter[*test](2, tt.length)

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
					t.Errorf("Filter.Contains() = %v, want %v", got, false)
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
