package flux

import (
	"e-highway-collector/lib/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"time"
)

// InfluxWriter InfluxDB写入
type InfluxWriter struct {
	writeApi api.WriteAPI
}

func (i *InfluxWriter) Write(msg Line) {
	time.Sleep(1 * time.Second) // separate points by 1 second
	logger.Info("record stored")
	i.writeApi.WriteRecord(string(msg))
	i.writeApi.Flush()
}

func MakeInfluxWriter(token string, url string, org string, bucket string) *InfluxWriter {
	client := influxdb2.NewClient(url, token)
	writeApi := client.WriteAPI(org, bucket)
	return &InfluxWriter{
		writeApi: writeApi,
	}
}
