// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package quake

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hpe-hcss/quake-client/pkg/terraform/configuration"
)

func TestGetConfigFromMeta(t *testing.T) {
	c := new(configuration.Config)
	m := map[string]interface{}{
		configuration.KeyForGLClientMap(): c,
	}
	tests := []struct {
		name   string
		input  interface{}
		output interface{}
		err    bool
	}{
		{
			name:   "success, *configuration.Config",
			input:  c,
			output: c,
			err:    false,
		},
		{
			name:   "success, map[configuration.KeyForGLClientMap()]*configuration.Config",
			input:  m,
			output: c,
			err:    false,
		},
		{
			name:   "error, not *configuration.Config",
			input:  struct{ name string }{name: "blah"},
			output: nil,
			err:    true,
		},
		{
			name: "error, not map[configuration.KeyForGLClientMap()]*configuration.Config",
			input: map[string]interface{}{
				"blah": c,
			},
			output: nil,
			err:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, e := getConfigFromMeta(tt.input)
			if tt.err {
				assert.NotNil(t, e)
				assert.Nil(t, o)
			} else {
				assert.Nil(t, e)
				assert.Equal(t, tt.output, o)
			}
		})
	}
}

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
