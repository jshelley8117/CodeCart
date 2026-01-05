package utils

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Env string // "dev" or "prod"
}

func NewLogger(cfg Config) (*zap.Logger, error) {
	env := strings.ToLower(strings.TrimSpace(cfg.Env))

	var zcfg zap.Config
	if env == "dev" || env == "development" || env == "" {
		zcfg = zap.NewDevelopmentConfig()
	} else {
		zcfg = zap.NewProductionConfig()
	}

	if lvl := strings.TrimSpace(os.Getenv("LOG_LEVEL")); lvl != "" {
		var parsed zapcore.Level
		if err := parsed.Set(lvl); err == nil {
			zcfg.Level = zap.NewAtomicLevelAt(parsed)
		}
	}

	return zcfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// uniquely defined key type - using a custom type prevents collisions with other packages that also store values in context
type ctxKey struct{}

// returns a new context derived from the original context, with the logger stored inside it.
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// looks up a logger in the context
// if found, return it
// if not found, return the fallback logger (the base logger injected into the struct)
func FromContext(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return fallback
}
