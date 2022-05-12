package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const Format = "2006-01-02 15:04:05.000"

const (
	LoginParams     = iota + 1 // 輸入資料有誤
	LoginMailFormat            // 信箱格式不正確
	LoginCreateFail            // 建立帳號失敗
	LoginRefresh               // 請先重新整理
	LoginRegisted              // 該信箱已經註冊
	LoginSuccess               // 建立帳號失敗
	ValidateNoToken            // token 遺失
	ValidateExpired            // token 過期
	ValidateFail               // 驗證失敗
	ValidateSuccess            // 驗證成功
	SendMail                   // 送出信件
)

func Init(name string) {
	//建立資料夾
	p := path.Join("log", name)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		errDir := os.MkdirAll(p, os.ModePerm)
		if errDir != nil {
			panic(errDir)
		}
	}

	//文字輸出
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: Format,
		ForceColors:     true,
		FullTimestamp:   true,
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	//記錄檔格式 log.2021_10_15-00_00_00.log
	fileName := path.Join(p, "log.%Y_%m_%d-%H_%M_%S.log")
	writer, err := rotatelogs.New(
		fileName,
		rotatelogs.WithMaxAge(7*24*time.Hour),       //保留10天前日誌
		rotatelogs.WithRotationTime(1*24*time.Hour), //每1天生成新的日誌
	)

	if err != nil {
		panic(err)
	}

	log := &logger{
		rotateLogs: writer,
		formatter: &logrus.JSONFormatter{
			TimestampFormat: Format,
		},
	}

	// 給寫檔
	logrus.AddHook(log)

	// gin.DefaultWriter = io.MultiWriter(log.rotateLogs)
}

type logger struct {
	rotateLogs *rotatelogs.RotateLogs
	formatter  *logrus.JSONFormatter
}

func (l *logger) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *logger) Fire(entry *logrus.Entry) error {
	if entry.Level == logrus.ErrorLevel {
		_, file, line, ok := runtime.Caller(6)
		if ok {
			entry.Data["line"] = fmt.Sprintf("%s:%d", file, line)
		}
	}

	buf, err := l.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = l.rotateLogs.Write(buf)
	return err
}
