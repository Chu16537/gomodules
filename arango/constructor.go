package arango

import (
	"context"
	"sync"
	"time"

	"github.com/Chu16537/gomodules/env"
	"github.com/Chu16537/gomodules/util"

	"github.com/Chu16537/gomodules/logger"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type manager struct {
	db     driver.Database
	ctx    context.Context
	config env.ArangoDB
}

var Instance *manager
var once sync.Once

func Init(c context.Context, conf env.ArangoDB) {
	once.Do(func() {
		logger.Debug("arango Init Start")
		Instance = &manager{
			ctx:    c,
			config: conf,
		}

		err := create()

		if err != nil {
			go util.Retry(Instance.ctx, create)
		}
	})
}

func Get() *manager {
	if Instance.db != nil {
		return Instance
	}
	return nil
}

// 連線實做
func create() error {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{Instance.config.Addr},
	})

	if err != nil {
		logger.Error("arango Init NewConnection Fail: %v", err)
		return err
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(Instance.config.Username, Instance.config.Password),
	})

	if err != nil {
		logger.Error("arango Init NewClient Fail: %v", err)
		return err
	}

	db, err := client.Database(Instance.ctx, Instance.config.Database)
	if err != nil {
		logger.Error("arango Init Database Fail: %v", err)
		return err
	}

	Instance.db = db

	logger.Info("arango Init Success")
	go pingLoop(client)
	return nil
}

// ping db 是否存活
func pingLoop(c driver.Client) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-Instance.ctx.Done():
			return
		case <-tick.C:
			if _, err := c.Version(Instance.ctx); err != nil {
				Instance.db = nil
				go util.Retry(Instance.ctx, create)
				return
			}
		}
	}
}
