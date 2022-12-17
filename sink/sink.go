package sink

import (
	"e-highway-collector/flux"
	"e-highway-collector/lib/logger"
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMQSink struct {
	TargetQueue    *amqp.Queue
	WorkingChannel *amqp.Channel
}

func MakeRabbitMQSink(name, url string) *RabbitMQSink {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic("conn create error")
	}
	ch, err := conn.Channel()
	q, _ := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil)
	return &RabbitMQSink{
		TargetQueue:    &q,
		WorkingChannel: ch,
	}
}

func (r *RabbitMQSink) Send(msg flux.Line) {
	msg.Timestamp = int(time.Now().Unix())
	mess, err := json.Marshal(msg)
	if err != nil {
		logger.Error("message encoding failed")
	}
	_ = r.WorkingChannel.Publish(
		"",
		r.TargetQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        mess,
		})
}
