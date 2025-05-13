package lslog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"runtime"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/go-ini/ini"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/petermattis/goid"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/xormplus/core"
)

// GetCurrentExecutePath 获取当时运行程序的目录 get current exe path
func getCurrentExecutePath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetExecuteFileName 获取当时运行程序的名称 get exe file name
func getExecuteFileName() (string, error) {
	dir := strings.Split(os.Args[0], string(os.PathSeparator))
	filename := dir[len(dir)-1]
	surffix := path.Ext(filename)
	return strings.TrimSuffix(filename, surffix), nil
}

// Logger a logger
var Logger = logrus.New()

// LinkseeXormLogger a logger to impelment xorm core.ILogger
type LinkseeXormLogger struct {
	log     *logrus.Logger
	level   core.LogLevel
	showSQL bool
}

// NewLinkseeXormLogger new a instance of MyXormLogger
func NewLinkseeXormLogger() *LinkseeXormLogger {
	var l core.LogLevel
	switch Logger.Level {
	case logrus.DebugLevel:
		l = core.LOG_DEBUG
	case logrus.InfoLevel:
		l = core.LOG_INFO
	case logrus.WarnLevel:
		l = core.LOG_WARNING
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		l = core.LOG_ERR
	default:
		l = core.LOG_UNKNOWN
	}
	return &LinkseeXormLogger{
		log:   Logger,
		level: l,
	}
}

// Debug to output debug logs
func (s *LinkseeXormLogger) Debug(v ...interface{}) {
	s.log.Debug(v...)
}

// Debugf to format & output debug logs
func (s *LinkseeXormLogger) Debugf(format string, v ...interface{}) {
	s.log.Debugf(format, v...)
}

// Error to output error logs
func (s *LinkseeXormLogger) Error(v ...interface{}) {
	s.log.Error(v...)
}

// Errorf to format & output error logs
func (s *LinkseeXormLogger) Errorf(format string, v ...interface{}) {
	s.log.Errorf(format, v...)
}

// Info to output info logs
func (s *LinkseeXormLogger) Info(v ...interface{}) {
	s.log.Info(v...)
}

// Infof to format & output info logs
func (s *LinkseeXormLogger) Infof(format string, v ...interface{}) {
	s.log.Infof(format, v...)
}

// Warn to output warning logs
func (s *LinkseeXormLogger) Warn(v ...interface{}) {
	s.log.Warn(v...)
}

// Warnf to format & output warning logs
func (s *LinkseeXormLogger) Warnf(format string, v ...interface{}) {
	s.log.Warnf(format, v...)
}

// Level implement core.ILogger
func (s *LinkseeXormLogger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *LinkseeXormLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

// ShowSQL implement core.ILogger
func (s *LinkseeXormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *LinkseeXormLogger) IsShowSQL() bool {
	return s.showSQL
}

// GetLogPrefix get log prefix string
func GetLogPrefix() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	format := fmt.Sprintf("[TID:%04d] [%12s:%04d] ", GetTid(), filename[0:strings.LastIndex(filename, ".go")], line)
	return format
}

// GetTid get gorouting id
func GetTid() int64 {
	return goid.Get()
}

// GetCodeLineNo get source code file an line no
func GetCodeLineNo(needext bool) string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	if needext {
		return fmt.Sprintf("[%12s:%04d]", filename, line)
	}
	return fmt.Sprintf("[%12s:%04d]", filename[0:strings.LastIndex(filename, ".go")], line)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	Logger.Out = writer
}

// New create a new instance of logger
func init() {
	level := "INFO"
	logFileName := ""
	logFormat := "2006/01/02 15:04:05"
	cfgfile, _ := getCurrentExecutePath()
	cfgfile = fmt.Sprintf("%s%clog.cfg", cfgfile, os.PathSeparator)
	iniFile, err := ini.InsensitiveLoad(cfgfile)
	if err == nil {
		level = strings.ToUpper(strings.TrimSpace(iniFile.Section("LOG").Key("LogLevel").String()))
		if level == "" {
			level = "INFO"
		}

		target := strings.ToUpper(strings.TrimSpace(iniFile.Section("LOG").Key("LogTarget").String()))
		if target == "FILE" {
			//	logFileName, _ = getExePath()
			logFileName, _ = getExecuteFileName()
		}

		format := strings.ToUpper(strings.TrimSpace(iniFile.Section("LOG").Key("LogTimeFormat").String()))
		if format != "" {
			logFormat = format
		}
	}

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	formatter.ForceColors = false
	formatter.TimestampFormat = logFormat
	if logFileName == "" {
		formatter.ForceColors = true
		// Set specific colors for prefix and timestamp
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle:    "blue+b",
			TimestampStyle: "white+h",
		})
	}
	Logger.Formatter = formatter
	// then wrap the log output with it
	Logger.Out = ansicolor.NewAnsiColorWriter(os.Stdout)
	switch strings.ToUpper(level) {
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	if logFileName != "" {
		//logPath, _ := getExePath()
		logPath, _ := getCurrentExecutePath()
		logPath = path.Join(logPath, "log")

		if exist := fileutil.IsExist(logPath); !exist {
			_ = os.MkdirAll(logPath, 0744)
		}
		logPath = path.Join(logPath, logFileName)
		var writer *rotatelogs.RotateLogs
		loc, err := time.LoadLocation("Asia/Shanghai") // 设置时区
		if err != nil {
			writer, _ = rotatelogs.New(
				logPath+"_log_%Y%m%d"+".txt",
				rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
			)
		} else {
			writer, _ = rotatelogs.New(
				logPath+"_log_%Y%m%d"+".txt",
				rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
				rotatelogs.WithLocation(loc),              // 当地时间
			)
		}
		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}, formatter)
		Logger.AddHook(lfHook)
	}
}
