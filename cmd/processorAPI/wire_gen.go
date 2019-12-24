// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/challenge/payment-processor/internal/app"
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/http"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/challenge/payment-processor/internal/pkg/processor/api"
	"github.com/challenge/payment-processor/internal/pkg/webserver"
)

// Injectors from wire.go:

func initializeConfig(path string) (config.Configuration, error) {
	configuration, err := config.NewConfig(path)
	if err != nil {
		return config.Configuration{}, err
	}
	return configuration, nil
}

func initializeAppender(path string) (logging.FileAppender, error) {
	configuration, err := config.NewConfig(path)
	if err != nil {
		return logging.FileAppender{}, err
	}
	fileAppender := logging.NewFileAppender(configuration)
	return fileAppender, nil
}

func initializeApp(path string) (*app.PaymentProcessorApp, error) {
	actorsMap := processor.NewActorsMap()
	acquirerActors := processor.NewAcquirerActors(actorsMap)
	configuration, err := config.NewConfig(path)
	if err != nil {
		return nil, err
	}
	loggingLevels := config.NewLoggingLevels(configuration)
	loggerProcessor := logging.NewLoggerProcessor(loggingLevels)
	paymentProcessorService := processor.NewPaymentProcessorService(acquirerActors, loggerProcessor)
	loggerAPI := logging.NewLoggerAPI(loggingLevels)
	controller := api.NewController(paymentProcessorService, loggerAPI)
	loggerWebServer := logging.NewWebServer(loggingLevels)
	paramWebServer := webserver.NewParamWebServer(controller, configuration, loggerWebServer)
	webServer := webserver.NewWebServer(paramWebServer)
	redisDB := config.NewRedisDB(configuration)
	loggerRedisDB := logging.NewRedisDB(loggingLevels)
	dbConnection, err := redis.NewDBConnection(redisDB, loggerRedisDB)
	if err != nil {
		return nil, err
	}
	loggerCardRepository := logging.NewLoggerCardRepository(loggingLevels)
	cardRepositoryRedis := processor.NewCardRepositoryRedis(dbConnection, redisDB, configuration, loggerCardRepository)
	loggerHTTP := logging.NewLoggerHTTP(loggingLevels)
	service := http.NewHTTPService(loggerHTTP)
	acquirerParameter := processor.NewAcquirerParameter(cardRepositoryRedis, service)
	stoneAcquirerWorkers := processor.NewStoneAcquirerWorkers(acquirerActors, acquirerParameter, configuration)
	cieloAcquirerWorkers := processor.NewCieloAcquirerWorkers(acquirerActors, acquirerParameter, configuration)
	logger := logging.NewLogger(loggingLevels)
	paymentProcessorApp := app.NewPaymentProcessorApp(webServer, stoneAcquirerWorkers, cieloAcquirerWorkers, logger)
	return paymentProcessorApp, nil
}
