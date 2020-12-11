package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func GetEnvVar(envVar, defaultValue string) string {
	value, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}

	return value
}

func read(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		errors.Wrap(err, fmt.Sprintf("couldn't read the file %s", path))
	}
	return content, nil
}

func ReadYAML(path string, dest interface{}) error {
	content, err := read(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(content, dest); err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't read the YAML file"))
	}
	return nil
}
