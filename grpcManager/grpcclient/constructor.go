package grpcclient

import (
	"context"

	"github.com/Chu16537/gomodules/env"
	"github.com/Chu16537/gomodules/logger"
	"github.com/Chu16537/gomodules/util"
	"google.golang.org/grpc"
)

type Handle struct {
	ctx    context.Context
	retry  *util.Retry
	config env.Grpc
	conn   *grpc.ClientConn
}

func (h *Handle) GetConn() *grpc.ClientConn {
	if h.conn == nil {
		return nil
	}

	return h.conn
}

func Create(c context.Context, config env.Grpc) (error, *Handle) {
	logger.Debug("grpc Create Start")
	defer logger.Debug("grpc Creat End")

	h := &Handle{
		ctx:    c,
		config: config,
	}

	h.retry = util.NewRetry(config.RetryCount, config.RetryTime)

	err, _ := h.retry.RunCount(func() (error, interface{}) {
		err := newClient(h)
		return err, nil
	})

	if err != nil {
		return err, nil
	}

	return nil, h
}

func newClient(h *Handle) error {
	dialOpt := grpc.WithInsecure()

	conn, err := grpc.Dial(h.config.Addr, dialOpt)

	if err != nil {
		return err
	}

	h.conn = conn
	return nil
}
