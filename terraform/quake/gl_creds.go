package quake

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Gljwt struct {
	SpaceName string `yaml:"space_name,omitempty"`
	ProjectID string `yaml:"project_id"`
	RestURL   string `yaml:"rest_url"`
	Token     string `yaml:"access_token"`
}

func loadGLConfig(dir string) (*Gljwt, error) {
	f, err := os.Open(filepath.Clean(filepath.Join(dir, ".gltform")))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseGLStream(f)
}

func parseGLStream(s io.Reader) (*Gljwt, error) {
	contents, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, err
	}

	q := &Gljwt{}
	if err = yaml.Unmarshal(contents, q); err != nil {
		return nil, err
	}
	return q, nil
}
