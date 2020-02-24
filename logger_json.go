package gin

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ptechen/encoding"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// TraceLevel defines trace log level.
	TraceLevel = iota - 1
	// DebugLevel defines debug log level.
	DebugLevel
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled
)

const (
	// TimeFormatUnix defines a time format that makes time fields to be
	// serialized as Unix timestamp integers.
	TimeFormatUnix = ""

	// TimeFormatUnixMs defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in milliseconds.
	TimeFormatUnixMs = "UNIXMS"

	// TimeFormatUnixMicro defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in microseconds.
	TimeFormatUnixMicro = "UNIXMICRO"
)

// JsonLoggerConfig defines the config for Logger middleware.
type JsonLoggerConfig struct {
	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string

	// IsConsole is whether to enable terminal printing.
	IsConsole bool

	// LogLevel is log level.
	LogLevel int8

	// Caller is whether to enable log tracking.
	Caller bool

	// LogColor is whether to enable log color.
	LogColor bool

	// LogWriteSize is sets the size of the log write pipeline.
	LogWriteSize int

	// LogTimeFieldFormat is a formatted layout of time fields in the log.
	LogTimeFieldFormat string

	// LogExpDays is number of days the log is kept.
	LogExpDays int64

	// LogLimitSize is the limit size of the log file, for example 1G and 512MB.
	LogLimitSize string

	logFilePath string

	logDir string

	logName string

	logLimitNums int64
}

// JsonLogger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func JsonLogger(conf ...JsonLoggerConfig) HandlerFunc {
	if len(conf) == 0 {
		return JsonLoggerWithConfig(JsonLoggerConfig{})
	}
	return JsonLoggerWithConfig(conf[0])
}

var once sync.Once
var logger *zerolog.Logger
var onceLog sync.Once

type TraceParams struct {
	StartTime time.Time
	Path      string
	ClientIp  string
	Method    string
}

// JsonLoggerWithConfig instance a Logger middleware with config.
func JsonLoggerWithConfig(conf JsonLoggerConfig) HandlerFunc {

	if conf.Output == nil {
		conf.Output = DefaultWriter
	}

	once.Do(func() {
		conf.initLogConfig()
		data, ok := conf.Output.(*os.File)
		if ok && conf.IsConsole == false {
			onceLog.Do(func() {
				conf.logFilePath = data.Name()
				data := strings.SplitAfter(data.Name(), "/")
				conf.logDir, conf.logName = data[0], data[1]
				conf.monitor()
			})
		}
	})

	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		params := &TraceParams{
			StartTime: start,
			Path:      c.Request.URL.String(),
			ClientIp:  c.ClientIP(),
			Method:    c.Request.Method,
		}

		traceId := createUuid(params)
		c.Logger = log.With().
			Str("trace_id", traceId).
			Str("path", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Int("body_size", c.Writer.Size()).
			Logger()

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				//isTerm:  isTerm,
				Keys: c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()
			if param.ErrorMessage == "" {
				c.Logger.Info().Dur("latency", param.Latency).
					Int("status", param.StatusCode).
					Interface("keys", c.Keys).Send()
			} else {
				c.Logger.Info().Dur("latency", param.Latency).
					Int("status", param.StatusCode).
					Err(errors.New(param.ErrorMessage)).
					Interface("keys", c.Keys).Send()
			}
		}
	}
}

func (p *JsonLoggerConfig) initLogConfig() {
	logger = &log.Logger
	p.setLogFileSize()
	p.setLogTimeFiledFormat()
	p.setLoglevel()
	p.setLogExpDays()
	p.setCaller()
	p.setLogWriteSize()
	p.setOutput()
}

func (p *JsonLoggerConfig) setCaller() {
	if p.Caller {
		*logger = logger.With().Caller().Logger()
	}
}

func (p *JsonLoggerConfig) setOutput() {
	var w io.Writer
	if p.IsConsole {
		if p.LogColor {
			w = diode.NewWriter(zerolog.ConsoleWriter{Out: DefaultWriter}, p.LogWriteSize,
				10*time.Millisecond, func(missed int) {
					logger.Warn().Msgf("Logger Dropped %d messages", missed)
				})

		} else {
			w = diode.NewWriter(DefaultWriter, p.LogWriteSize, 10*time.Millisecond, func(missed int) {
				logger.Warn().Msgf("Logger Dropped %d messages", missed)
			})
		}

	} else {
		w = diode.NewWriter(p.Output, p.LogWriteSize, 10*time.Millisecond, func(missed int) {
			logger.Warn().Msgf("Logger Dropped %d messages", missed)
		})
	}

	*logger = logger.Output(w)
}

func (p *JsonLoggerConfig) setLoglevel() {
	if p.LogLevel < -1 && p.LogLevel > 7 {
		p.LogLevel = 0
	}
	*logger = logger.Level(zerolog.Level(p.LogLevel))
}

func (p *JsonLoggerConfig) setLogWriteSize() {
	if p.LogWriteSize == 0 {
		p.LogWriteSize = 1000
	}
}

func (p *JsonLoggerConfig) setLogFileSize() {
	if p.LogLimitSize == "" {
		p.LogLimitSize = "1G"
	}
	if strings.Contains(p.LogLimitSize, "G") {
		n, _ := strconv.Atoi(strings.Split(p.LogLimitSize, "G")[0])
		p.logLimitNums = int64(n) * 1024 * 1024 * 1024

	} else if strings.Contains(p.LogLimitSize, "MB") {
		n, _ := strconv.Atoi(strings.Split(p.LogLimitSize, "MB")[0])
		p.logLimitNums = int64(n) * 1024 * 1024
	} else {
		panic("please input ")
	}
}

func (p *JsonLoggerConfig) setLogTimeFiledFormat() {
	zerolog.TimeFieldFormat = p.LogTimeFieldFormat
}

func (p *JsonLoggerConfig) setLogExpDays() {
	if p.LogExpDays == 0 {
		p.LogExpDays = 30
	}
}

func (p *JsonLoggerConfig) monitor() {
	t := time.NewTicker(time.Second * 3)
	del := time.NewTicker(time.Hour * 24)

	go func() {
		defer t.Stop()
		defer del.Stop()
		for {
			select {

			case <-t.C:
				isExist := p.isExist()
				if !isExist {
					p.setOutput()
				}
				size := p.checkFileSize()
				if size > p.logLimitNums {
					logger.Info().Msg("rename log file")
					p.rename2File()
					p.setOutput()
				}

			case <-del.C:
				p.deleteLogFile()
			}
		}
	}()
}

func (p *JsonLoggerConfig) isExist() bool {
	_, err := os.Stat(p.logFilePath)
	return err == nil || os.IsExist(err)
}

func (p *JsonLoggerConfig) checkFileSize() int64 {
	f, e := os.Stat(p.logFilePath)
	if e != nil {
		return 0
	}
	return f.Size()
}

func (p *JsonLoggerConfig) rename2File() {
	now := time.Now()
	newLogFileName := fmt.Sprintf("%s.%s", p.logFilePath, now.Format("2006-01-02 15:04:05"))
	_ = os.Rename(p.logFilePath, newLogFileName)
}

func (p *JsonLoggerConfig) deleteLogFile() {
	files, _ := ioutil.ReadDir(p.logDir)
	for _, file := range files {
		if file.IsDir() {
			// DO
		} else {
			if file.Name() != p.logName && strings.Contains(file.Name(), p.logName) {
				createTime := strings.Split(file.Name(), p.logName+".")[1]
				date, err := time.Parse("2006-01-02 15:04:05", createTime)
				if err != nil {
					logger.Err(err).Msg("log file time format err")
					continue
				}
				dateUnix := date.Unix()
				currentUnix := time.Now().Unix()
				if currentUnix-dateUnix > p.LogExpDays*60*60*24 {
					currentFileName := p.logDir + "/" + file.Name()
					err = os.Remove(currentFileName)
					if err != nil {
						logger.Err(err).Msgf("remove %s failed", currentFileName)
					}
					logger.Info().Msgf("remove %s success", currentFileName)
				}
			}
		}
	}
}

func createUuid(params interface{}) (uuidStr string) {
	data, err := encoding.JSON.Marshal(params)
	if err != nil {
		return uuidStr
	}
	uuidStr = uuid.NewMD5(uuid.UUID{}, data).String()
	return uuidStr
}
