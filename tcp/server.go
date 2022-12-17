package tcp

import (
	"context"
	"e-highway-collector/config"
	"e-highway-collector/core/scheduler"
	"e-highway-collector/core/worker"
	"e-highway-collector/flux"
	"e-highway-collector/interface/tcp"
	"e-highway-collector/lib/logger"
	"e-highway-collector/sink"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(
	cfg *Config,
	handler tcp.Handler) error {
	closeChan := make(chan struct{})

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// 监听 Signal 管道信号的协程
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {

	// 监听 Close 管道信号的协程
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()

	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	var waitDone sync.WaitGroup

	ctx := context.Background()
	ch := make(chan flux.Line, 100)
	influxWriter := flux.MakeInfluxWriter(
		config.Properties.InfluxToken,
		config.Properties.InfluxUrl,
		config.Properties.InfluxOrg,
		config.Properties.InfluxBucket)

	rabbitMQSink := sink.MakeRabbitMQSink(
		config.Properties.RabbitMQQueueName,
		config.Properties.RabbitMQUrl)

	concurrentScheduler := scheduler.MakeConcurrentScheduler(
		ch,
		config.Properties.WorkerNum,
		worker.Worker,
		func(currentWorker, workerCount int) int {
			currentWorker++
			return currentWorker % workerCount
		},
		influxWriter,
		rabbitMQSink)

	go concurrentScheduler.Run()

	// 网络层主循环，向通信管道中写入
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		logger.Info("accepted link")
		waitDone.Add(1)

		// 每监听到一个新的连接，创建一个新的协程处理该客户端
		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn, ch)
		}()
	}
	waitDone.Wait()
}
