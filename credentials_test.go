package yc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparable(t *testing.T) {
	for _, tt := range []struct {
		lhs   interface{}
		rhs   interface{}
		equal bool
	}{
		{
			lhs:   mustNewClient(),
			rhs:   mustNewClient(),
			equal: true,
		},
		{
			lhs:   mustNewClient(WithEndpoint("test")),
			rhs:   mustNewClient(WithEndpoint("test")),
			equal: true,
		},
		{
			lhs:   mustNewClient(WithEndpoint("test")),
			rhs:   mustNewClient(),
			equal: false,
		},
		{
			lhs:   NewInstanceServiceAccount(),
			rhs:   mustNewClient(),
			equal: false,
		},
	} {
		t.Run("", func(t *testing.T) {
			if tt.equal {
				assert.EqualValues(t, tt.lhs, tt.rhs)
			} else {
				assert.NotEqualValues(t, tt.lhs, tt.rhs)
			}
		})
	}
}
