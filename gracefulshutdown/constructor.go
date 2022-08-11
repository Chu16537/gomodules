package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"

	"github.com/Chu16537/gomodules/logger"
)

var ctx context.Context
var cancel context.CancelFunc
var shutdownChan chan os.Signal

func Init() {
	logger.Debug("gracefulshutdown Init")

	ctx, cancel = context.WithCancel(context.Background())

	shutdownChan = make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	go waitShutdown()
}

func GetContext() context.Context {
	return ctx
}

func Shutdown() {
	shutdownChan <- os.Interrupt
}

func waitShutdown() {
	<-shutdownChan
	signal.Stop(shutdownChan)

	logger.Debug("gracefulshutdown waitShutdown")

	defer cancel()
}
