package grpcManager

import (
	"context"
	"net"
	"time"

	"github.com/Chu16537/gomodules/env"
	"github.com/Chu16537/gomodules/logger"
	"github.com/Chu16537/gomodules/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Handle struct {
	ctx    context.Context
	retry  *util.Retry
	config env.Grpc
	lis    net.Listener
	server *grpc.Server
}

func (h *Handle) GetCtx() context.Context {
	return h.ctx
}

func (h *Handle) GetServer() *grpc.Server {
	return h.server
}

func Create(c context.Context, config env.Grpc) (error, *Handle) {
	logger.Debug("grpc Create Start")
	defer logger.Debug("grpc Creat End")

	h := &Handle{
		ctx:    c,
		config: config,
	}

	h.retry = util.NewRetry(config.RetryCount, config.RetryTime)

	err, _ := h.retry.Run(h.ctx, func() (error, interface{}) {
		err := newServer(h)
		return err, nil
	})

	if err != nil {
		return err, nil
	}

	return nil, h
}

// 創建 grcp
func newServer(h *Handle) error {
	lis, err := net.Listen("tcp", h.config.Addr)

	if err != nil {
		logger.Error("newServer net.Listen to Fail: %v", err)
		return err
	}

	h.lis = lis

	h.server = grpc.NewServer()
	reflection.Register(h.server)

	go checkLoop(h)

	logger.Debug("grpc newServer Success")

	return nil
}

// 檢查連線是否存在
func checkLoop(h *Handle) {
	tick := time.NewTicker(h.config.RetryTime * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-tick.C:
			if err := h.server.Serve(h.lis); err != nil {

				select {
				case <-h.ctx.Done():
					return
				default:
				}

				logger.Error("failed to server: %v", err)

				retryErr, _ := h.retry.Run(h.ctx, func() (error, interface{}) {
					h.lis.Close()
					h.server.Stop()
					err := newServer(h)
					return err, nil
				})

				if retryErr == nil {
					return
				}
			}

		}
	}
}
