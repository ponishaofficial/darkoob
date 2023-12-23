package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitFunc(t *testing.T) {
	splitFunc, ok := FuncMaps["split"]
	assert.True(t, ok)
	split := splitFunc.(func(s any, sep string, i int) string)

	type test struct {
		input any
		sep   string
		index int
		want  string
	}

	tests := []test{
		{"text without separator", "@", 0, "text without separator"},
		{"text-with-multiple-separator", "-", 0, "text"},
		{"my@awesome.email", "@", 0, "my"},
		{"my@awesome.email", "@", 1, "awesome.email"},
		{"my@awesome.email", "@", 2, "my"},
		{12345, "@", 2, ""},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Equalf(t, tc.want, split(tc.input, tc.sep, tc.index), "split(%s, %s, %d)", tc.input, tc.sep, tc.input)
		})
	}
}

func TestRandIntFunc(t *testing.T) {
	randIntFunc, ok := FuncMaps["randInt"]
	assert.True(t, ok)
	randInt := randIntFunc.(func(n ...int) []int)

	type test struct {
		input []int
	}

	tests := []test{
		{[]int{}},
		{[]int{1, 1, 1}},
		{[]int{1, 2, 3, 2, 1}},
	}

	lessThan := func(a, b []int) bool {
		if len(a) != len(b) {
			return false
		}

		for i := range a {
			if b[i] >= a[i] || b[i] < 0 {
				return false
			}
		}

		return true
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Truef(t, lessThan(tc.input, randInt(tc.input...)), "randInt(%+v)", tc.input)
		})
	}
}

func TestJoinFunc(t *testing.T) {
	joinFunc, ok := FuncMaps["join"]
	assert.True(t, ok)
	join := joinFunc.(func(sep string, s []any) string)

	type test struct {
		input []any
		sep   string
		want  string
	}

	tests := []test{
		{[]any{}, "@", ""},
		{[]any{1, 2, 3}, ",", "1,2,3"},
		{[]any{"1"}, ",", "1"},
		{[]any{"1", 2}, ",", "1,2"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Equalf(t, tc.want, join(tc.sep, tc.input), "join(%s, %+v)", tc.sep, tc.input)
		})
	}
}
