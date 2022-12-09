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

func (i *InfluxWriter) Write(msg [][]byte) {
	var con string
	for _, line := range msg {
		con += string(line)
	}
	logger.Info(con)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"gate": "1",
		}
		fields := map[string]interface{}{
			"count": value,
		}
		point := write.NewPoint("example", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second
		logger.Info("record stored")

		if err := i.writeApi.WritePoint(context.Background(), point); err != nil {
			logger.Error(err.Error())
		}
	}
}

func MakeInfluxWriter(token string, url string, org string, bucket string) *InfluxWriter {
	client := influxdb2.NewClient(url, token)
	writeApi := client.WriteAPIBlocking(org, bucket)
	return &InfluxWriter{
		writeApi: writeApi,
	}
}
