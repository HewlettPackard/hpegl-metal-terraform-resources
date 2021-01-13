package configuration

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const QjwtExtension = ".qjwt"

// Qjwt declares the contents of the login file.
type Qjwt struct {
	RestURL     string `yaml:"rest_url"`
	OriginalURL string `yaml:"original_url"`
	User        string `yaml:"user"`
	Token       string `yaml:"jwt"`
	MemberID    string `yaml:"member_id"`
	NoTLS       bool   `yaml:"no_tls"`
}

func loadConfig(dir string) (*Qjwt, error) {
	f, err := os.Open(filepath.Clean(filepath.Join(dir, QjwtExtension)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseStream(f)
}

func parseStream(s io.Reader) (*Qjwt, error) {
	contents, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, err
	}

	q := &Qjwt{}
	if err = yaml.Unmarshal(contents, q); err != nil {
		return nil, err
	}
	return q, nil
}
