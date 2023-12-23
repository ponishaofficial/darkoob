package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_processLine(t *testing.T) {
	type testCase[T interface{ ~string | any }] struct {
		data map[string]map[string]any
		line T
		want T
	}
	tests := []testCase[string]{
		{nil, `{{ randInt "," 1 1 1 }}`, "0,0,0"},
		{nil, `{{ split "my,name,is" "," 1 }}`, "name"},
		{nil, `{{ join "," 1 2 3 4 5 }}`, "1,2,3,4,5"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Equalf(t, tt.want, processLine(tt.data, tt.line), "processLine(%v, %v)", tt.data, tt.line)
		})
	}
}
