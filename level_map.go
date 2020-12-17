package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// AtomicLevelMap string level mapping zap AtomicLevel
	AtomicLevelMap = map[string]zap.AtomicLevel{
		"debug":  zap.NewAtomicLevelAt(zap.DebugLevel),
		"info":   zap.NewAtomicLevelAt(zap.InfoLevel),
		"warn":   zap.NewAtomicLevelAt(zap.WarnLevel),
		"error":  zap.NewAtomicLevelAt(zap.ErrorLevel),
		"dpanic": zap.NewAtomicLevelAt(zap.DPanicLevel),
		"panic":  zap.NewAtomicLevelAt(zap.PanicLevel),
		"fatal":  zap.NewAtomicLevelAt(zap.FatalLevel),
	}

	// ZapcoreLevelMap string level mapping zapcore.Level
	ZapcoreLevelMap = map[string]zapcore.Level{
		"debug":  zap.DebugLevel,
		"info":   zap.InfoLevel,
		"warn":   zap.WarnLevel,
		"error":  zap.ErrorLevel,
		"dpanic": zap.DPanicLevel,
		"panic":  zap.PanicLevel,
		"fatal":  zap.FatalLevel,
	}
)
