package log

import (
    "validation_service/pkg/config"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
    var (
        instance  *zap.Logger
        logConfig zap.Config
    )
    if config.Settings.Env == "production" {
        logConfig = zap.NewProductionConfig()
    } else {
        logConfig = zap.NewDevelopmentConfig()
    }
    logConfig.EncoderConfig.TimeKey = "timestamp"
    logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

    instance, err := logConfig.Build()
    if err != nil {
        panic(err)
    }

    defer instance.Sync()

    Logger = instance.Sugar()
}
