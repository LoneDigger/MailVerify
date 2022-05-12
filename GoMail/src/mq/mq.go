package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"me.mail/src/config"
)

// 訊息類型
const (
	SendMail = "send_mail"
)

const (
	exchangeName = "amq.topic"
	queueName    = "reg.mail"
)

// 郵件
type Mail struct {
	Ip         string    `json:"ip"`
	Mail       string    `json:"mail"`
	Token      string    `json:"token"`
	UserID     int       `json:"user_id"`
	CreateTime time.Time `json:"createtime"`
}

type Connect struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	msgCh chan Mail
}

func NewConnect(msg config.Msg, msgCh chan Mail) *Connect {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/reg", msg.Username, msg.Password, msg.Host, msg.Port))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Connect{
		conn:  conn,
		ch:    ch,
		msgCh: msgCh,
	}
}

func (conn *Connect) loop(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var msg Mail
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			logrus.WithError(err).Error()
		} else {
			conn.msgCh <- msg
		}
	}
}

func (conn *Connect) Send(obj interface{}) error {
	b, _ := json.Marshal(obj)

	return conn.ch.Publish(exchangeName, queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(b),
	})
}

func (conn *Connect) Start() {
	// 建立佇列
	_, err := conn.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	{
		err = conn.ch.QueueBind(
			queueName,
			queueName,
			exchangeName,
			false,
			nil,
		)

		if err != nil {
			panic(err)
		}
	}

	msgs, err := conn.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		panic(err)
	}

	go conn.loop(msgs)
}
