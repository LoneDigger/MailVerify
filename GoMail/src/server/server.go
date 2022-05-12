package server

import (
	"github.com/sirupsen/logrus"
	"me.mail/src/config"
	"me.mail/src/dialer"
	"me.mail/src/logger"
	"me.mail/src/mq"
)

type Server struct {
	dialer *dialer.Dialer
	conn   *mq.Connect
	msgCh  chan mq.Mail
}

func NewServer(cfg config.Config) *Server {
	s := new(Server)
	s.msgCh = make(chan mq.Mail)

	s.conn = mq.NewConnect(cfg.Msg, s.msgCh)
	s.dialer = dialer.NewMail(cfg.Host, cfg.Mail)

	return s
}

func (s *Server) Start() {
	s.conn.Start()
	go s.loop()
}

func (s *Server) loop() {
	for msg := range s.msgCh {
		fields := logrus.Fields{
			"ip":          msg.Ip,
			"mail":        msg.Mail,
			"create_time": msg.CreateTime,
			"action":      logger.SendMail,
		}

		err := s.dialer.Send(msg)
		if err != nil {
			logrus.WithError(err).WithFields(fields).Error()
		} else {
			logrus.WithFields(fields).Info()
		}
	}
}
