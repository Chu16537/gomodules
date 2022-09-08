package env

import "time"

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
