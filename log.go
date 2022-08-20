package function

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger(exch, symbol string) *zap.SugaredLogger {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: true,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:      "N",
			CallerKey:    "C",
			MessageKey:   "M",
			LineEnding:   zapcore.DefaultLineEnding,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"e": exch, "s": symbol,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
