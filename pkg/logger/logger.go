package logger

import (
	"errors"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"time"
)

var (
	Level  = logrus.DebugLevel
	Logger *logrus.Logger
	Writer io.Writer
)

func InitLog(level logrus.Level) error {
	logPath := viper.GetString("log_path")
	logName := viper.GetString("log_name")
	if logName == "" {
		logName = "server.log"
	}
	if logPath == "" {
		logPath = "./logs"
	}
	Logger = logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	//formatter.SetColorScheme(&prefixed.ColorScheme{
	//	PrefixStyle:    "white+h",
	//	InfoLevelStyle: "white",
	//	TimestampStyle: "black+h"})
	Logger.SetFormatter(formatter)
	Logger.SetOutput(os.Stderr)
	Logger.SetLevel(level)

	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			return errors.New("Failed to create log folder")
		}
	}
	baseLogPath := path.Join(logPath, logName)
	Writer, err = rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		return err
	}

	fileformatter := new(prefixed.TextFormatter)
	fileformatter.FullTimestamp = true
	fileformatter.TimestampFormat = "2006-01-02 15:04:05"
	fileformatter.DisableColors = true

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: Writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  Writer,
		logrus.WarnLevel:  Writer,
		logrus.ErrorLevel: Writer,
		logrus.FatalLevel: Writer,
		logrus.PanicLevel: Writer,
	}, fileformatter)
	Logger.AddHook(lfHook)
	return nil
}
