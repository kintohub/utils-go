package klog

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "VERBOSE":
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "DEBUG":
		// TODO: Decide if we want this because its not json but it highlights logs much better!
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msgf("zerolog logger initialized with %s settings", zerolog.GlobalLevel())
}

func Verbose(msg string) {
	log.Log().Msg(msg)
}

func Verbosef(format string, args ...interface{}) {
	log.Log().Msgf(format, args...)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func WarnfWithErr(err error, format string, args ...interface{}) {
	log.Warn().Err(err).Msgf(format, args...)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func ErrorfWithErr(err error, format string, args ...interface{}) {
	log.Error().Err(err).Msgf(format, args...)
}

func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

func FatalfWithErr(err error, format string, args ...interface{}) {
	log.Fatal().Err(err).Msgf(format, args...)
}

func Panic(msg string) {
	log.Panic().Msg(msg)
}

func Panicf(format string, args ...interface{}) {
	log.Panic().Msgf(format, args...)
}

func PanicfWithError(err error, format string, args ...interface{}) {
	log.Panic().Err(err).Msgf(format, args...)
}

// function used to measure the time a function took to execute
// to measure a function time call this function as the first time with the following:
// Ex: `defer LogDuration(time.Now(), "doSomethingFunc")` will output: "doSomethingFunc took 82us"
func LogDuration(start time.Time, name string) {
	elapsed := time.Since(start)
	Debugf("%s took %s", name, elapsed)
}

func GetLogger() zerolog.Logger {
	return log.Logger
}
