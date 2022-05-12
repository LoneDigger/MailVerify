package server

import (
	"github.com/gin-gonic/gin"
	"me.user/src/config"
	"me.user/src/db"
	"me.user/src/mq"
	"me.user/src/redis"
)

// 回應
type RegisterResult struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type Server struct {
	e *gin.Engine
	d *db.DB

	r  *redis.Register
	v  *redis.Validation
	rq *mq.Connect
}

func NewServer(cfg config.Config) *Server {
	s := new(Server)

	s.d = db.NewDB(cfg.Postgres)
	s.r = redis.NewRegister(cfg.Register)
	s.v = redis.NewValidation(cfg.Validate)

	s.rq = mq.NewConnect(cfg.Msg)

	s.e = gin.Default()
	s.e.LoadHTMLGlob("templates/*")

	s.e.GET("/create", s.createView)
	s.e.POST("/create", s.create)

	s.e.GET("/register/:token", s.registerView)
	s.e.GET("/register", s.registerView)

	return s
}

func (s *Server) Start() error {
	s.rq.Start()

	return s.e.Run(":80")
}
