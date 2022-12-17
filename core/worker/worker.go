package worker

import (
	"e-highway-collector/flux"
	"e-highway-collector/lib/logger"
	"e-highway-collector/sink"
	"fmt"
)

func Worker(influxWriter *flux.InfluxWriter, sinkClient *sink.RabbitMQSink, line flux.Line) error {
	go func() {
		sinkClient.Send(line)
		influxWriter.Write(line)
		logContent := fmt.Sprintf("point(measurement=%s, tags=%v, fields=%v, timestamp=%d)",
			line.Measurement, line.Tags, line.Fields, line.Timestamp)
		logger.Info("->" + logContent)
	}()
	return nil
}
