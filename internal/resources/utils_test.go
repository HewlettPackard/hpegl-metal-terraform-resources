// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifference(t *testing.T) {
	tests := []struct {
		name   string
		slice1 []string
		slice2 []string
		retval []string
	}{
		{
			name:   "Test1ExpectDiff",
			slice1: []string{"aaa", "bbb", "ccc"},
			slice2: []string{"aaa", "bbb"},
			retval: []string{"ccc"},
		},
		{
			name:   "Test2NoDiff",
			slice1: []string{"aaa", "bbb"},
			slice2: []string{"aaa", "bbb", "cccc"},
			retval: []string{},
		},
		{
			name:   "Test3Same",
			slice1: []string{"aaa", "bbb"},
			slice2: []string{"aaa", "bbb"},
			retval: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret := difference(tt.slice1, tt.slice2)
			assert.Equal(t, tt.retval, ret)
		})
	}
}
