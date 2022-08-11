package arango

import (
	"context"
	"sync"
	"time"

	"github.com/Chu16537/gomodules/gracefulshutdown"
	"github.com/Chu16537/gomodules/logger"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type manager struct {
	db  driver.Database
	ctx context.Context
}

var Instance *manager
var once sync.Once

func Init() {
	once.Do(func() {
		logger.Debug("arango Init Start")
		Instance = &manager{
			ctx: gracefulshutdown.GetContext(),
		}

		err := create()

		if err != nil {
			go retry()
		}
	})
}

// 連線實做
func create() error {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{conf.Env.ArangoDB.Addr},
	})

	if err != nil {
		logger.Error("arango Init NewConnection Fail: %v", err)
		return err
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(conf.Env.ArangoDB.Username, conf.Env.ArangoDB.Password),
	})

	if err != nil {
		logger.Error("arango Init NewClient Fail: %v", err)
		return err
	}

	db, err := client.Database(Instance.ctx, conf.Env.ArangoDB.Database)
	if err != nil {
		logger.Error("arango Init Database Fail: %v", err)
		return err
	}

	Instance.db = db

	logger.Info("arango Init Success")
	go pingLoop(client)
	return nil
}

// 重新連線
func retry() {

	nowCount := 0
	tick := time.NewTicker(conf.Env.ArangoDB.RetryTime * time.Millisecond)
	defer tick.Stop()

	for nowCount < conf.Env.ArangoDB.RetryCount {
		select {
		case <-Instance.ctx.Done():
			gracefulshutdown.Shutdown()
			return
		case <-tick.C:
			err := create()
			if err == nil {
				return
			}
			nowCount++
		}

	}

	gracefulshutdown.Shutdown()
}

// ping db 是否存活
func pingLoop(c driver.Client) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-Instance.ctx.Done():
			gracefulshutdown.Shutdown()
			return
		case <-tick.C:
			if _, err := c.Version(Instance.ctx); err != nil {
				go retry()
				return
			}
		}
	}
}
