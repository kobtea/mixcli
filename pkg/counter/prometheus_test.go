package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountExprs(t *testing.T) {
	tests := []struct {
		exprs  []string
		result Counter
	}{
		{
			[]string{
				"foo",
				"bar",
				"foo",
			},
			Counter{
				"foo": 2,
				"bar": 1,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, CountExprs(test.exprs), test.result)
	}
}

func TestCountSelectors(t *testing.T) {
	tests := []struct {
		exprs  []string
		result Counter
	}{
		{
			[]string{
				`foo{n="1"} + bar{n="2"}`,
				"foo",
			},
			Counter{
				`foo{n="1"}`: 1,
				`bar{n="2"}`: 1,
				"foo":        1,
			},
		},
	}
	for _, test := range tests {
		res, err := CountSelectors(test.exprs)
		assert.Nil(t, err)
		assert.Equal(t, res, test.result)
	}
}

func TestCountMetricNames(t *testing.T) {
	tests := []struct {
		exprs      []string
		resMetrics Counter
		resNoNames Counter
	}{
		{
			[]string{
				`foo{n="1"} + bar{n="2"} > 0`,
				`foo + {__name__="baz"}`,
				`{__name__=~"amb.+"}`,
			},
			Counter{
				"foo": 2,
				"bar": 1,
				"baz": 1,
			},
			Counter{
				`{__name__=~"amb.+"}`: 1,
			},
		},
	}
	for _, test := range tests {
		resMetrics, resNoNames, err := CountMetricNames(test.exprs)
		assert.Nil(t, err)
		assert.Equal(t, resMetrics, test.resMetrics)
		assert.Equal(t, resNoNames, test.resNoNames)
	}
}
