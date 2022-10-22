// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"
	usecase "nextclan/transaction-gateway/transaction-validate-service/internal/usecase/transaction"
	"nextclan/transaction-gateway/transaction-validate-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, vt usecase.ReceiveRawTransactionFromQueue) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	//how well is the http server running
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	//
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
