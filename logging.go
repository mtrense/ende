package ende

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	initialized bool = false
	logger      zerolog.Logger
)

func L() zerolog.Logger {
	if !initialized {
		level, err := zerolog.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			level = zerolog.WarnLevel
		}
		zerolog.SetGlobalLevel(level)
		output := viper.GetString("logfile")
		var writer io.Writer
		if output == "-" {
			writer = zerolog.ConsoleWriter{Out: os.Stderr}
		} else {
			writer, err = os.OpenFile(output, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
			if err != nil {
				panic(err)
			}
		}
		logger = zerolog.New(writer).With().Timestamp().Logger()
	}
	return logger
}

