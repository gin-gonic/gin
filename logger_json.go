package gin

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
	"time"
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
}

// ZeroLogger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func JsonLogger(conf ... JsonLoggerConfig) HandlerFunc {
	if len(conf) == 0 {
		return JsonLoggerWithConfig(JsonLoggerConfig{})
	}
	return JsonLoggerWithConfig(conf[0])
}

var once sync.Once
var logger *zerolog.Logger

// JsonLoggerWithConfig instance a Logger middleware with config.
func JsonLoggerWithConfig(conf JsonLoggerConfig) HandlerFunc {

	if conf.Output == nil {
		conf.Output = DefaultWriter
	}

	once.Do(func() {
		logger = &log.Logger
		conf.loglevel()
		conf.caller()
		conf.logWriteSize()
		conf.output()
	})

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		c.Logger = log.With().
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

func (p *JsonLoggerConfig) caller() {
	if p.Caller {
		*logger = logger.With().Caller().Logger()
	}
}

func (p *JsonLoggerConfig) output() {
	var w io.Writer
	if p.IsConsole {

		if p.LogColor {
			w = diode.NewWriter(zerolog.ConsoleWriter{Out: p.Output}, p.LogWriteSize,
			10*time.Millisecond, func(missed int) {
				logger.Warn().Msgf("Logger Dropped %d messages", missed)
			})

		} else {
			w = diode.NewWriter(p.Output, p.LogWriteSize, 10*time.Millisecond, func(missed int) {
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

func (p *JsonLoggerConfig) loglevel() {
	if p.LogLevel < -1 && p.LogLevel > 7 {
		p.LogLevel = 0
	}
	*logger = logger.Level(zerolog.Level(p.LogLevel))
}

func (p *JsonLoggerConfig) logWriteSize() {
	if p.LogWriteSize == 0 {
		p.LogWriteSize = 1000
	}
}