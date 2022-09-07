package env

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 讀取
func Get(path string, conf interface{}) error {

	configByte, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configByte, conf)

	if err != nil {
		return err
	}

	return nil
}
