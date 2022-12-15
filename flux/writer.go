package flux

import (
	"context"
	"e-highway-collector/lib/logger"
	"fmt"
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
	time.Sleep(1 * time.Second) // separate points by 1 second
	point := write.NewPoint(
		msg.Measurement,
		msg.Tags,
		msg.Fields,
		time.Now())
	fmt.Printf("point(measurement=%s, tags=%v, fields=%v, timestamp=%d)",
		msg.Measurement, msg.Tags, msg.Fields, msg.Timestamp)
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
