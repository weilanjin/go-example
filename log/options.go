package logger

import (
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	FMT_CONSOLE    = "console"
	FMT_JSON       = "json"
	OUTPUT_STD     = "stdout"
	OUTPUT_STD_ERR = "stderr"
)

type Options struct {
	OutputPath      []string `mapstructure:"output-path"`
	ErrorOutputPath []string `mapstructure:"error-output-path"`
	Level           string   `mapstructure:"level"`
	Format          string   `mapstructure:"format"`
	Name            string   `mapstructure:"name"`
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	options := Options{
		Level:           zapcore.InfoLevel.String(),
		Format:          FMT_CONSOLE,
		OutputPath:      []string{OUTPUT_STD},
		ErrorOutputPath: []string{OUTPUT_STD_ERR},
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &options
}

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func (o *Options) Validate() []error {
	var errs []error
	logFmt := strings.ToLower(o.Format)
	if logFmt != FMT_CONSOLE && logFmt != FMT_JSON {
		errs = append(errs, fmt.Errorf("invalid log format: %s", o.Format))
	}
	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.OutputPath, "log-output-path", o.OutputPath, "log output path")
	fs.StringSliceVar(&o.ErrorOutputPath, "log-error-output-path", o.ErrorOutputPath, "log error output path")
	fs.StringVar(&o.Level, "log-level", o.Level, "log level")
	fs.StringVar(&o.Format, "log-format", o.Format, "log format")
	fs.StringVar(&o.Name, "log-name", o.Name, "log name")
}