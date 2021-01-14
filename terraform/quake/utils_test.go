package quake

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/quattronetworks/quake-client/pkg/terraform/configuration"
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
