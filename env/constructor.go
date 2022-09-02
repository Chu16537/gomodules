package env

import (
	"io/ioutil"
	"time"

	"github.com/Chu16537/gomodules/logger"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ArangoDB *ArangoDB `yaml:"Arango,omitempty"`
	Grpc     *Grpc     `yaml:"Grpc,omitempty"`
}

// ArangoDB DB
type ArangoDB struct {
	Addr         string        `yaml:"addr,omitempty"`
	Database     string        `yaml:"database,omitempty"`
	Username     string        `yaml:"username,omitempty"`
	Password     string        `yaml:"password,omitempty"`
	HttpProtocol string        `yaml:"httpProtocol,omitempty"`
	RetryCount   int           `yaml:"retryCount,omitempty"`
	RetryTime    time.Duration `yaml:"retryTime,omitempty"`
}

type Grpc struct {
	Addr       string        `yaml:"addr,omitempty"`
	RetryCount int           `yaml:"retryCount,omitempty"`
	RetryTime  time.Duration `yaml:"retryTime,omitempty"`
}

// 讀取
func Load() (*Config, error) {
	configByte, err := ioutil.ReadFile("env.yaml")

	if err != nil {
		logger.Error("env Init Err:%v", err)
		configByte, err = ioutil.ReadFile("env.template.yaml")
	}

	if err != nil {
		logger.Error("env Init Err:%v", err)
		return nil, err
	}

	env := &Config{}
	err = yaml.Unmarshal(configByte, env)

	if err != nil {
		return nil, err
	}

	return env, nil
}
