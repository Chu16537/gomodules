package arango

import (
	"context"
	"time"

	"github.com/Chu16537/gomodules/env"
	"github.com/Chu16537/gomodules/util"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"

	"github.com/Chu16537/gomodules/logger"
)

type Handle struct {
	ctx    context.Context
	retry  *util.Retry
	config env.ArangoDB
	db     driver.Database
}

func Create(c context.Context, config env.ArangoDB) (error, *Handle) {
	logger.Debug("arango Creat Start")
	defer logger.Debug("arango Creat End")

	h := &Handle{
		ctx:    c,
		config: config,
	}

	h.retry = util.NewRetry(config.RetryCount, config.RetryTime)

	err, _ := h.retry.Run(h.ctx, func() (error, interface{}) {
		err := connect(h)
		return err, nil
	})

	if err != nil {
		return err, nil
	}

	return nil, h
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

	go h.checkLoop(client)
	h.db = db

	logger.Debug("arango connect Success")
	return nil
}

// ping db 是否存活
func (h *Handle) checkLoop(client driver.Client) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-tick.C:
			if _, err := client.Version(h.ctx); err != nil {
				logger.Error("arango ping Fail")

				retryErr, _ := h.retry.Run(h.ctx, func() (error, interface{}) {
					err := connect(h)
					return err, nil
				})

				if retryErr == nil {
					return
				}

			}
		}
	}
}
