package function

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func setupLogger(exch, symbol string, sig *Signal) *zap.SugaredLogger {
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
			EncodeCaller: LineNumCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"e": exch, "s": symbol,
			"m": sig.Strategy, "a": sig.Action,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}

func LineNumCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	p := caller.TrimmedPath()
	idx := strings.IndexByte(p, '/')
	if idx != -1 {
		p = p[idx+1:]
	}
	enc.AppendString(p)
}
