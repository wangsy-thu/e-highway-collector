package flux

import "e-highway-collector/lib/logger"

func Worker(influxWriter *InfluxWriter, line Line) error {
	logger.Info("receive new line")
	influxWriter.Write(line)
	return nil
}
