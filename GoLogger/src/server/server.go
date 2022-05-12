package server

import (
	"github.com/sirupsen/logrus"
	"me.log/src/config"
	"me.log/src/db"
	"me.log/src/mq"
)

type Server struct {
	conn  *mq.Connect
	db    *db.MongoDB
	msgCh chan db.Box
}

func NewServer(cfg config.Config) *Server {
	s := new(Server)
	s.msgCh = make(chan db.Box)
	s.db = db.NewMongoDB(cfg.Mongo)
	s.conn = mq.NewConnect(cfg.Msg, s.msgCh)

	return s
}

func (s *Server) Start() {
	s.conn.Start()
	go s.loop()
}

func (s *Server) loop() {
	for obj := range s.msgCh {
		err := s.db.Write(obj)
		if err != nil {
			logrus.WithError(err).WithField("box", obj).Error()
		} else {
			logrus.WithField("box", obj).Info()
		}
	}
}
