package main

import (
	"e-highway-collector/flux"
	"e-highway-collector/lib/logger"
)

func main() {
	token := "VVAwQEWXAVOzp0eXkkFK-0aho4qrBWFK3wCN88XwzuXLvZfhtGZJmizPZg76GwxEJIxlGy4KvQonvQWbVlST7w=="
	url := "http://127.0.0.1:8086"
	bucket := "example"
	org := "neu-thu"
	influxWriter := flux.MakeInfluxWriter(token, url, org, bucket)
	content := [][]byte{[]byte("hello"), []byte("world")}
	influxWriter.Write(content)
	logger.Info("Write Finished")
}
