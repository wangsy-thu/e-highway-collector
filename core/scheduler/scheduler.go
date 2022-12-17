package scheduler

import (
	flux2 "e-highway-collector/flux"
	_ "e-highway-collector/interface/flux"
	"e-highway-collector/lib/logger"
	"e-highway-collector/sink"
)

type Saver func(influxWriter *flux2.InfluxWriter, sink *sink.RabbitMQSink, line flux2.Line) error
type ScheduleFunc func(currentWorker, workerCount int) int

type ConcurrentScheduler struct {
	WorkerCount int
	WorkingChan map[int]chan flux2.Line
	Worker      Saver
	NextWorker  ScheduleFunc
	SourceChan  chan flux2.Line
	FluxClient  *flux2.InfluxWriter
	SinkClient  *sink.RabbitMQSink
}

func MakeConcurrentScheduler(source chan flux2.Line,
	workerCount int,
	worker Saver,
	nextWorker ScheduleFunc,
	fluxClient *flux2.InfluxWriter,
	sinkClient *sink.RabbitMQSink) *ConcurrentScheduler {
	cs := &ConcurrentScheduler{
		WorkerCount: workerCount,
		Worker:      worker,
		NextWorker:  nextWorker,
		SourceChan:  source,
		FluxClient:  fluxClient,
		SinkClient:  sinkClient,
	}
	return cs
}

func (scheduler *ConcurrentScheduler) Run() {
	// 创建工作协程组
	scheduler.WorkingChan = make(map[int]chan flux2.Line, scheduler.WorkerCount)
	for i := 0; i < scheduler.WorkerCount; i++ {
		scheduler.WorkingChan[i] = make(chan flux2.Line)
	}
	// 启动工作协程组
	for i := 0; i < scheduler.WorkerCount; i++ {
		go func(workerIdx int) {
			for {
				payload := <-scheduler.WorkingChan[workerIdx]
				err := scheduler.Worker(scheduler.FluxClient, scheduler.SinkClient, payload)
				if err != nil {
					logger.Error("Unknown Error")
				}
			}
		}(i)
	}

	// 确定下一个工作协程
	cw := 0
	for {
		nw := scheduler.NextWorker(cw, scheduler.WorkerCount)
		cw = nw
		scheduler.WorkingChan[nw] <- <-scheduler.SourceChan
	}
}
