package cardcaldav

import (
	"github.com/charmbracelet/log"
	"github.com/rs/zerolog"
	log0 "github.com/rs/zerolog/log"
	"strings"
)

type zeroToCharmLogger struct{ logger *log.Logger }

func (z *zeroToCharmLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if len(s) <= 3 {
		return len(p), nil
	}
	inLevel := s[:3]
	var level log.Level
	switch inLevel {
	case "TRC", "DBG":
		level = log.DebugLevel
	case "INF":
		level = log.InfoLevel
	case "WRN":
		level = log.WarnLevel
	case "ERR":
		level = log.ErrorLevel
	case "FTL", "PNC":
		level = log.FatalLevel
	}
	z.logger.Helper()
	translator := z.logger.With()
	translator.SetCallerFormatter(func(s string, i int, s2 string) string {
		return "tokidoki internal"
	})
	translator.Log(level, strings.TrimSpace(s[4:]))
	return len(p), nil
}

func SetupLogger(logger *log.Logger) {
	log0.Logger = log0.Output(zerolog.ConsoleWriter{
		Out: &zeroToCharmLogger{logger},
		FormatTimestamp: func(i interface{}) string {
			return ""
		},
	})
}
