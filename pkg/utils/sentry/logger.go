package sentry

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sentryToZerologLevel = map[sentry.Level]zerolog.Level{
	sentry.LevelDebug:   zerolog.DebugLevel,
	sentry.LevelInfo:    zerolog.InfoLevel,
	sentry.LevelWarning: zerolog.WarnLevel,
	sentry.LevelError:   zerolog.ErrorLevel,
	sentry.LevelFatal:   zerolog.FatalLevel,
}

// Event is a top-level function for dealing with errors in a robust manner. It logs
// the error using zerolog at the specified level, sends an error to Sentry if the hub
// is available, adds the error to the gin context if it's available and performs other
// tasks related to monitoring and alerting of errors in the Ensign project.
//
// This should only be used if the error needs to generate an alert; otherwise use
// zerolog directly rather than using this Event type.
//
// The sentry level is mapped to the zerolog level, which means the zerolog.TraceLevel
// and zerolog.PanicLevel are not available in this Event type.
//
// Not safe for concurrent use!
type Event struct {
	zero  *zerolog.Event
	extra map[string]interface{}
	level sentry.Level
	hub   *sentry.Hub
	ginc  *gin.Context
	err   *ServiceError
}

// The ctx should be either a gin.Context or a context.Context; the hub is extracted
// from the context if it was set by middleware or interceptors. Once the event is
// created it must be sent using Msg or Msgf.
func CreateEvent(level sentry.Level, ctx interface{}) *Event {
	event := &Event{
		zero:  log.WithLevel(sentryToZerologLevel[level]),
		extra: make(map[string]interface{}),
		level: level,
	}

	// Attempt to fetch the hub from the context
	switch c := ctx.(type) {
	case *gin.Context:
		event.hub = sentrygin.GetHubFromContext(c)
		event.ginc = c
	case context.Context:
		event.hub = sentry.GetHubFromContext(c)
	case *sentry.Hub:
		event.hub = c
	case nil:
		event.hub = sentry.CurrentHub().Clone()
	}

	return event
}

// Reports a debug level event to Sentry and logs a debug message. Use this method when
// the debug message should produce an alert that the team can take action on (which
// should happen only very rarely in code). Most of the time you should use zerolog.Debug
// directly unless this is at the top level of the stack.
func Debug(ctx interface{}) *Event {
	return CreateEvent(sentry.LevelDebug, ctx)
}

// Reports an info level event to Sentry and logs an info message. Use this method when
// the info message should produce an alert that the team can take action on (which
// should happen very rarely in code and is probably related to a third party service
// such as rate limits or usage thresholds). Most of the time you should use zerolog.Info
// directly unless this is at the top level of the stack.
func Info(ctx interface{}) *Event {
	return CreateEvent(sentry.LevelInfo, ctx)
}

// Report a warning level event to Sentry and logs a warning messages. Use this method
// on top level service handlers to produce alerts that something is going wrong in the
// code such as bad requests or not found errors. The team will likely not take action
// on these errors but will get a general sense of what is going on in code. When not
// in a service handler it is better to use zerolog.Warn directly.
func Warn(ctx interface{}) *Event {
	return CreateEvent(sentry.LevelWarning, ctx)
}

// Report an error to Sentry and log an error message. This is the most commonly used
// method for Sentry on top level service handlers and is intended to produce alerts
// that something is going wrong and that the team needs to handle it. When not in a
// service handler, feel free to use zerolog.Error but probably zerolog.Warn is more
// appropriate for most cases.
func Error(ctx interface{}) *Event {
	return CreateEvent(sentry.LevelError, ctx)
}

// Report a critical error to Sentry and log a fatal error message. While this method
// will not cause the process to exit, it should create a serious alert that will cause
// on call personnel to immediately act. Use with care!
func Fatal(ctx interface{}) *Event {
	return CreateEvent(sentry.LevelFatal, ctx)
}

func (e *Event) Err(err error) *Event {
	if err != nil {
		e.err = &ServiceError{err: err}
	}
	e.zero = e.zero.Err(err)
	return e
}

func (e *Event) Errs(errs []error) *Event {
	if len(errs) > 0 {
		e.err = &ServiceError{err: multierror.Append(nil, errs...)}
	}
	e.zero = e.zero.Errs("errors", errs)
	return e
}

func (e *Event) Str(key, value string) *Event {
	e.extra[key] = value
	e.zero = e.zero.Str(key, value)
	return e
}

func (e *Event) Int(key string, value int) *Event {
	e.extra[key] = value
	e.zero = e.zero.Int(key, value)
	return e
}

func (e *Event) Uint8(key string, value uint8) *Event {
	e.extra[key] = value
	e.zero = e.zero.Uint8(key, value)
	return e
}

func (e *Event) Int32(key string, value int32) *Event {
	e.extra[key] = value
	e.zero = e.zero.Int32(key, value)
	return e
}

func (e *Event) Int64(key string, value int64) *Event {
	e.extra[key] = value
	e.zero = e.zero.Int64(key, value)
	return e
}

func (e *Event) Uint64(key string, value uint64) *Event {
	e.extra[key] = value
	e.zero = e.zero.Uint64(key, value)
	return e
}

func (e *Event) ULID(key string, value ulid.ULID) *Event {
	s := value.String()
	e.extra[key] = s
	e.zero = e.zero.Str(key, s)
	return e
}

func (e *Event) Bytes(key string, value []byte) *Event {
	e.extra[key] = base64.RawURLEncoding.EncodeToString(value)
	e.zero = e.zero.Bytes(key, value)
	return e
}

func (e *Event) Bool(key string, value bool) *Event {
	e.extra[key] = value
	e.zero = e.zero.Bool(key, value)
	return e
}

func (e *Event) Dict(key string, dict *Dictionary) *Event {
	for dkey, dval := range dict.extra {
		e.extra[fmt.Sprintf("%s_%s", key, dkey)] = dval
	}
	e.zero = e.zero.Dict(key, dict.dict)
	return e
}

// Finalizes the event and sends it to Sentry and Zerolog
func (e *Event) Msg(msg string) {
	// Update the error with the context message
	if e.err != nil {
		e.err.msg = msg

		// If a gin context is available set the error on it
		if e.ginc != nil {
			e.ginc.Error(e.err)
		}
	}

	// If a hub is available send the message to sentry.
	if e.hub != nil {
		e.hub.ConfigureScope(func(scope *sentry.Scope) {
			if len(e.extra) > 0 {
				scope.SetContext("error", e.extra)
			}
			scope.SetLevel(e.level)
		})

		if e.err != nil {
			e.hub.CaptureException(e.err)
		} else {
			e.hub.CaptureMessage(msg)
		}
	}

	// Log the message to zerolog
	e.zero.Msg(msg)
}

// Finalizes the event with the format string and arguments then sends it to Sentry and Zerolog.
func (e *Event) Msgf(format string, args ...interface{}) {
	// Update the error with the context message
	if e.err != nil {
		e.err.msg = format
		e.err.args = args

		// If a gin context is available set the error on it
		if e.ginc != nil {
			e.ginc.Error(e.err)
		}
	}

	// If a hub is available send the message to sentry.
	if e.hub != nil {
		e.hub.ConfigureScope(func(scope *sentry.Scope) {
			if len(e.extra) > 0 {
				scope.SetContext("error", e.extra)
			}
			scope.SetLevel(e.level)
		})

		if e.err != nil {
			e.hub.CaptureException(e.err)
		} else {
			e.hub.CaptureMessage(fmt.Sprintf(format, args...))
		}
	}

	// Log the message to zerolog
	e.zero.Msgf(format, args...)
}

func With(ctx interface{}) *Logger {
	return &Logger{
		zero:  log.With(),
		extra: make(map[string]interface{}),
		ctx:   ctx,
	}
}

// Logger is an intermediate struct that holds a zero-log context.
type Logger struct {
	zero  zerolog.Context
	extra map[string]interface{}
	ctx   interface{}
}

func (l Logger) createEvent(level sentry.Level, zero *zerolog.Event) *Event {
	event := &Event{
		zero:  zero,
		extra: l.extra,
		level: level,
	}

	// Attempt to fetch the hub from the context
	switch c := l.ctx.(type) {
	case *gin.Context:
		event.hub = sentrygin.GetHubFromContext(c)
		event.ginc = c
	case context.Context:
		event.hub = sentry.GetHubFromContext(c)
	case *sentry.Hub:
		event.hub = c
	case nil:
		event.hub = sentry.CurrentHub().Clone()
	}

	return event
}

func (l Logger) Trace() *zerolog.Event {
	lg := l.zero.Logger()
	return lg.Trace()
}

func (l Logger) Debug() *Event {
	lg := l.zero.Logger()
	return l.createEvent(sentry.LevelDebug, lg.Debug())
}

func (l Logger) Info() *Event {
	lg := l.zero.Logger()
	return l.createEvent(sentry.LevelInfo, lg.Info())
}

func (l Logger) Warn() *Event {
	lg := l.zero.Logger()
	return l.createEvent(sentry.LevelWarning, lg.Warn())
}

func (l Logger) Error() *Event {
	lg := l.zero.Logger()
	return l.createEvent(sentry.LevelError, lg.Error())
}

func (l *Logger) Str(key, value string) *Logger {
	l.extra[key] = value
	l.zero = l.zero.Str(key, value)
	return l
}

func (l *Logger) Int(key string, value int) *Logger {
	l.extra[key] = value
	l.zero = l.zero.Int(key, value)
	return l
}

func (l *Logger) Bool(key string, value bool) *Logger {
	l.extra[key] = value
	l.zero = l.zero.Bool(key, value)
	return l
}

func (l *Logger) Dict(key string, dict *Dictionary) *Logger {
	for dkey, dval := range dict.extra {
		l.extra[fmt.Sprintf("%s_%s", key, dkey)] = dval
	}
	l.zero = l.zero.Dict(key, dict.dict)
	return l
}

func Dict() *Dictionary {
	return &Dictionary{
		dict:  zerolog.Dict(),
		extra: make(map[string]interface{}),
	}
}

type Dictionary struct {
	dict  *zerolog.Event
	extra map[string]interface{}
}

func (d *Dictionary) Str(key, val string) *Dictionary {
	d.extra[key] = val
	d.dict = d.dict.Str(key, val)
	return d
}

func (d *Dictionary) Int(key string, value int) *Dictionary {
	d.extra[key] = value
	d.dict = d.dict.Int(key, value)
	return d
}

func (d *Dictionary) Uint64(key string, value uint64) *Dictionary {
	d.extra[key] = value
	d.dict = d.dict.Uint64(key, value)
	return d
}
