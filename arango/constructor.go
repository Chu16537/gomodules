package arango

import (
	"context"
	"time"

	"github.com/Chu16537/gomodules/env"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"

	"github.com/Chu16537/gomodules/logger"
)

type Handle struct {
	db     driver.Database
	ctx    context.Context
	config env.ArangoDB
}

func Create(c context.Context, config env.ArangoDB) *Handle {
	logger.Debug("arango Creat Start")
	defer logger.Debug("arango Creat Success")

	h := &Handle{
		ctx:    c,
		config: config,
	}

	err := connect(h)

	if err != nil {
		go retry(h)
	}

	return h
}

// 連線實做
func connect(h *Handle) error {
	logger.Debug("arango connect Start")
	conn, connErr := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{h.config.Addr},
	})

	if connErr != nil {
		logger.Error("arango Init NewConnection Fail: %v", connErr)
		return connErr
	}

	client, clientErr := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(h.config.Username, h.config.Password),
	})

	if clientErr != nil {
		logger.Error("arango Init NewClient Fail: %v", clientErr)
		return clientErr
	}

	db, dbErr := client.Database(h.ctx, h.config.Database)

	if dbErr != nil {
		logger.Error("arango Init Database Fail: %v", dbErr)
		return dbErr
	}

	go pingLoop(h, client)
	h.db = db

	logger.Debug("arango connect Success")
	return nil
}

// ping db 是否存活
func pingLoop(h *Handle, client driver.Client) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-tick.C:
			if _, err := client.Version(h.ctx); err != nil {
				logger.Error("arango pingLoop Fail go retry")
				go retry(h)
				return
			}
		}
	}
}

// 重新創db
func retry(h *Handle) {
	tick := time.NewTicker(h.config.RetryTime * time.Millisecond)
	defer tick.Stop()

	select {
	case <-h.ctx.Done():
		return
	case <-tick.C:
		err := connect(h)
		if err == nil {
			return
		}
	}
}
