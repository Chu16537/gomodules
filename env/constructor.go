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
	// Redis     *Redis     `yaml:"redis,omitempty"`
	// Nats      *Nats      `yaml:"nats,omitempty"`

	RetryCount int           `yaml:"retryCount,omitempty"`
	RetryTime  time.Duration `yaml:"retryTime,omitempty"`
}

// ArangoDB DB
type ArangoDB struct {
	Addr         string `yaml:"addr,omitempty"`
	Database     string `yaml:"database,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	HttpProtocol string `yaml:"httpProtocol,omitempty"`
}

type Grpc struct {
	Addr string `yaml:"addr,omitempty"`
}

var Env *Config

// 讀取
func Init() error {
	configByte, err := ioutil.ReadFile("env.yaml")

	if err != nil {
		logger.Error("env Init Err:%v", err)
		configByte, err = ioutil.ReadFile("env.template.yaml")
	}

	if err != nil {
		logger.Error("env Init Err:%v", err)
		return err
	}

	err = yaml.Unmarshal(configByte, &Env)

	if err != nil {
		return err
	}

	return nil
}
