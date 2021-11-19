package log

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	InitDev()
}

func InitProd() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("logger failed to create: %v", err))
	}
	Logger = l.Sugar()
}

func InitDev() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("logger failed to create: %v", err))
	}
	Logger = l.Sugar()
}
