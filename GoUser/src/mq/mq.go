package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"me.user/src/config"
)

// 訊息類型
const (
	SendMail    = "send_mail"
	ValidateUrl = "validated_url"
)

const (
	exchangeName = "amq.topic"
	queueMail    = "reg.mail"
	queueLog     = "reg.log"
)

// 郵件
type Mail struct {
	Ip         string    `json:"ip"`
	Mail       string    `json:"mail"`
	Token      string    `json:"token"`
	UserID     int       `json:"user_id"`
	CreateTime time.Time `json:"createtime"`
}

// 驗證結果
type Validate struct {
	Ip         string    `json:"ip"`
	UserID     int       `json:"usesr_id"`
	CreateTime time.Time `json:"create_time"`
}

type Connect struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConnect(msg config.Msg) *Connect {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/reg", msg.Username, msg.Password, msg.Host, msg.Port))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Connect{
		conn: conn,
		ch:   ch,
	}
}

func (conn *Connect) send(t, queue string, obj interface{}) error {
	b, _ := json.Marshal(obj)

	return conn.ch.Publish(exchangeName, queue, false, false, amqp.Publishing{
		Type:        t,
		ContentType: "application/json",
		Body:        []byte(b),
	})
}

func (conn *Connect) SendLog(obj interface{}) error {
	return conn.send(ValidateUrl, queueLog, obj)
}

func (conn *Connect) SendMail(obj interface{}) error {
	return conn.send(SendMail, queueMail, obj)
}

func (conn *Connect) Start() {
}
