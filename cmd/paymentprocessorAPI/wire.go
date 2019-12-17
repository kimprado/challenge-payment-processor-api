// +build wireinject

package main

import (
	app "github.com/challenge/payment-processor/internal/app"
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/google/wire"
)

// initializeConfig inicializa Configuration
func initializeConfig(path string) (config config.Configuration, err error) {
	panic(wire.Build(app.AppSet))
}

// initializeAppender inicializa FileAppender
func initializeAppender(path string) (appender logging.FileAppender, err error) {
	panic(wire.Build(app.AppSet))
}

// initializeApp inicializa ExchangeApp
func initializeApp(path string) (a *app.PaymentProcessorApp, err error) {
	panic(wire.Build(app.AppSet))
}
