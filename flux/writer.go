package flux

import (
	"context"
	"e-highway-collector/lib/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

// InfluxWriter InfluxDB写入
type InfluxWriter struct {
	writeApi api.WriteAPIBlocking
}

func (i *InfluxWriter) Write(msg Line) {
	point := write.NewPoint(
		msg.Measurement,
		msg.Tags,
		msg.Fields,
		time.Now())
	err := i.writeApi.WritePoint(context.Background(), point)
	if err != nil {
		logger.Error("write error")
	}
}

func MakeInfluxWriter(token string, url string, org string, bucket string) *InfluxWriter {
	client := influxdb2.NewClient(url, token)
	writeApi := client.WriteAPIBlocking(org, bucket)
	return &InfluxWriter{
		writeApi: writeApi,
	}
}
