package server

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"me.user/src/db"
	"me.user/src/logger"
	"me.user/src/mq"
)

const (
	passwordFormat = "#%s#_+_+_<%s>"
	mailFormat     = "^[\\w\\.]{1,32}@[a-zA-Z0-9]{1,20}\\.[a-zA-Z0-9]{1,20}(.[a-zA-Z0-9]{1,20})?$"
)

// 正規表達式確認信相格式
func (s *Server) chekcMail(mail string) bool {
	return regexp.MustCompile(mailFormat).MatchString(mail)
}

// 建立帳號
func (s *Server) create(ctx *gin.Context) {
	var request db.User
	var result RegisterResult

	fields := logrus.Fields{
		"ip":   ctx.ClientIP(),
		"view": "create",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		result.Message = "輸入資料有誤"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginParams
		logrus.WithError(err).WithFields(fields).Error()
		return
	}

	fields["mail"] = request.Mail
	fields["name"] = request.Name

	// 確認郵件格式
	if !s.chekcMail(request.Mail) {
		result.Message = "信箱格式不正確"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginMailFormat
		logrus.WithFields(fields).Warn()
		return
	}

	// 密碼加密
	temp := fmt.Sprintf(passwordFormat, request.Mail, request.Password)
	pw, err := bcrypt.GenerateFromPassword([]byte(temp), bcrypt.DefaultCost)
	if err != nil {
		result.Message = "輸入資料有誤"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginParams
		logrus.WithError(err).WithFields(fields).Error()
		return
	}

	request.Hash = string(pw)

	// 確認 token
	b, err := s.r.Check(request.Token)
	if err != nil {
		result.Message = "建立帳號失敗"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginCreateFail
		logrus.WithFields(fields).Warn()
	}

	if !b {
		result.Message = "請先重新整理"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginRefresh
		logrus.WithFields(fields).Warn()
		return
	}

	// 確認郵件有無註冊過
	if s.d.CheckMail(request.Mail) {
		result.Message = "該信箱已經註冊"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginRegisted
		logrus.WithFields(fields).Warn()
		return
	}

	err = s.d.Create(&request)
	if err != nil {
		result.Message = "建立帳號失敗"
		ctx.JSON(http.StatusOK, result)

		fields["action"] = logger.LoginCreateFail
		logrus.WithError(err).WithFields(fields).Error()
	} else {
		fields["ip"] = request.UserID
		s.r.Remove(request.Token)

		err = s.pushMail(request.UserID, ctx.ClientIP(), request.Mail)
		if err != nil {
			result.Message = "建立帳號失敗"
			ctx.JSON(http.StatusOK, result)

			fields["action"] = logger.LoginCreateFail
			logrus.WithError(err).WithFields(fields).Error()
		} else {
			result.Message = "建立成功"
			result.Result = true
			ctx.JSON(http.StatusOK, result)

			fields["action"] = logger.LoginSuccess
			logrus.WithFields(fields).Info()
		}
	}
}

// 建立帳號畫面
func (s *Server) createView(ctx *gin.Context) {
	token := s.csrf("-")

	err := s.r.Add(token)
	if err != nil {
		logrus.WithError(err).Error()
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "create.html", gin.H{
		"token": token,
	})
}

// 雜湊碼
func (s *Server) csrf(arg string) string {
	t := fmt.Sprintf("%s#%s#", arg, time.Now().Format(logger.Format))
	return fmt.Sprintf("%x", md5.Sum([]byte(t)))
}

// 寄信驗證 郵件、紀錄
func (s *Server) pushMail(id int, ip, mail string) error {
	token := s.csrf(fmt.Sprintf("%d_%s", id, mail))

	err := s.v.Add(id, token)
	if err != nil {
		return err
	}

	return s.rq.SendMail(mq.Mail{
		Ip:         ip,
		Mail:       mail,
		Token:      token,
		UserID:     id,
		CreateTime: time.Now(),
	})
}

// 驗證連結紀錄
func (s *Server) pushValidate(id int, ip string) error {
	return s.rq.SendLog(mq.Validate{
		Ip:         ip,
		UserID:     id,
		CreateTime: time.Now(),
	})
}

// 驗證信箱畫面
func (s *Server) registerView(ctx *gin.Context) {
	token := ctx.Param("token")

	result := make(gin.H)
	fields := logrus.Fields{
		"ip":   ctx.ClientIP(),
		"view": "validate",
	}

	if token == "" {
		result["flag"] = false
		result["title"] = "驗證失敗"
		result["context"] = "遺失身份"

		fields["action"] = logger.ValidateNoToken
		logrus.WithFields(fields).Warn()
	} else {
		id, err := s.v.Check(token)

		if err != nil {
			result["flag"] = false
			result["title"] = "驗證失敗"
			result["context"] = "token 過期"

			if errors.Is(err, redis.Nil) {
				fields["action"] = logger.ValidateNoToken
				logrus.WithFields(fields).Warn()
			} else {
				fields["action"] = logger.ValidateFail
				logrus.WithError(err).WithFields(fields).Error()
			}
		} else {
			fields["id"] = id
			err = s.d.Update(id)

			if err != nil {
				result["flag"] = false
				result["title"] = "驗證失敗"
				result["context"] = "找不到註冊帳戶"

				fields["action"] = logger.ValidateFail
				logrus.WithFields(fields).Warn()
			} else {
				result["flag"] = true
				result["title"] = "驗證成功"
				result["context"] = ""

				fields["action"] = logger.ValidateFail
				logrus.WithFields(fields).Info()
				s.v.Remove(token)

				s.pushValidate(id, ctx.ClientIP())
			}
		}

		ctx.HTML(http.StatusOK, "validation.html", result)
	}
}
