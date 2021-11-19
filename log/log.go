package log

import (
	"fmt"
	"github.com/Velocityofpie/chaudr/config"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	var l *zap.Logger
	var err error
	if config.DebugMode {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("logger failed to create: %v", err))
	}

	Logger = l.Sugar()
}
