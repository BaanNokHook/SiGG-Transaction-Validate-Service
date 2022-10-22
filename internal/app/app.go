package app

import (
	"fmt"
	"nextclan/transaction-gateway/transaction-validate-service/config"
	v1 "nextclan/transaction-gateway/transaction-validate-service/internal/controller/http/v1"
	usecase "nextclan/transaction-gateway/transaction-validate-service/internal/usecase/transaction"
	"nextclan/transaction-gateway/transaction-validate-service/pkg/httpserver"
	"nextclan/transaction-gateway/transaction-validate-service/pkg/logger"
	"nextclan/transaction-gateway/transaction-validate-service/pkg/pusher"
	messaging "nextclan/transaction-gateway/transaction-validate-service/pkg/rabbitmq"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type sampleMessage struct {
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	fmt.Println("Starting App...")

	// Use case
	receiveValidatedTransactionUseCase := usecase.NewReceiveRawTransaction(l)
	// HTTP Server
	httpServer := initializeHttp(l, receiveValidatedTransactionUseCase, cfg)
	//Init client
	initializeMessaging(cfg, receiveValidatedTransactionUseCase)
	//Init pusher beam
	initialPusherBeam(cfg)
	// Shutdown
	ShutdownApplicationHandler(l, httpServer)
}

func initializeHttp(l *logger.Logger, vt *usecase.ReceiveRawTransactionFromQueueUseCase, cfg *config.Config) *httpserver.Server {
	handler := gin.New()
	v1.NewRouter(handler, l, vt)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	return httpServer
}

func initialPusherBeam(cfg *config.Config) {
	usecase.PushNotificationClient = pusher.NewPusherClient(cfg.PusherBeam.InstanceId, cfg.PusherBeam.SecretKey)
}

func initializeMessaging(cfg *config.Config, vt *usecase.ReceiveRawTransactionFromQueueUseCase) {
	//TODO dependency injection for usecase scope
	usecase.MessagingClient = &messaging.MessagingClient{}
	usecase.MessagingClient.Connect(cfg.RMQ.URL)
	usecase.MessagingClient.SubscribeToQueue("txt.gw", "topic", "raw.transaction", "transaction.validation.service", vt.Handle)
}

func ShutdownApplicationHandler(l *logger.Logger, httpServer *httpserver.Server) {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	}

	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = usecase.MessagingClient.Close()
}
