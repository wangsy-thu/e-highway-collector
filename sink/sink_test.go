package sink

import (
	"e-highway-collector/flux"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func TestSink(t *testing.T) {
	conn, err := amqp.Dial("amqp://WangY:wsy20010418@192.168.44.100")
	if err != nil {
		panic("conn create error")
	}
	tags := map[string]string{
		"tag1": "t1",
		"tag2": "t2",
	}
	fields := map[string]interface{}{
		"speed": 40,
		"plate": "辽A3631E",
	}
	l := flux.Line{
		Measurement: "hello",
		Tags:        tags,
		Fields:      fields,
		Timestamp:   123445,
	}
	result, _ := json.Marshal(l)
	fmt.Println("after json encoding: \n" + string(result))
	ch, err := conn.Channel()
	q, _ := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)
	_ = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        result,
		})
}
