package main

import (
	"fmt"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("logger failed to create: %v", err))
	}
	logger = l.Sugar()
}
